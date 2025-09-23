import time, json, psycopg2, threading
from psycopg2.extras import RealDictCursor
from concurrent.futures import ThreadPoolExecutor
from tqdm import tqdm
from playwright.sync_api import sync_playwright
from bs4 import BeautifulSoup

# ==============================================================
# Scraper
# ==============================================================

class CourseScraperPool:
    def __init__(self):
        self.local = threading.local()

    def get_page(self):
        if not hasattr(self.local, 'page'):
            pw = sync_playwright().start()
            browser = pw.chromium.launch(headless=True)
            self.local.playwright, self.local.browser, self.local.page = pw, browser, browser.new_page()
            self.local.page.goto("https://www1.reg.cmu.ac.th/registrationoffice/searchcourse.php")
        return self.local.page

    def cleanup(self):
        if hasattr(self.local, 'browser'):
            self.local.browser.close()
            self.local.playwright.stop()

    def scrape_course(self, code):
        try:
            p = self.get_page()
            p.fill("#fcourse", code); p.click("#button2")
            p.wait_for_load_state("networkidle", timeout=60000)
            soup = BeautifulSoup(p.content(), "html.parser")
            rows = soup.select(".tblCourse tr")
            data = []
            for r in rows:
                cols = r.find_all("td")
                if len(cols) < 12 or not cols[3].get_text(strip=True).isdigit():
                    continue
                rooms   = [x.get_text(strip=True) for x in cols[9].find_all(["div","span"]) if x.get_text(strip=True) not in ["","-"]]
                days    = [x.get_text(strip=True) for x in cols[7].find_all(["div","span"]) if x.get_text(strip=True) not in ["","-"]]
                times   = [x.get_text(strip=True) for x in cols[8].find_all(["div","span"]) if x.get_text(strip=True) not in ["","-"]]
                lecs    = list(cols[10].stripped_strings)
                for i,d in enumerate(days):
                    data.append(dict(
                        course_code=code, title=cols[2].get_text(strip=True),
                        lec_section=cols[3].get_text(strip=True) or None,
                        lab_section=cols[4].get_text(strip=True) or None,
                        credit=float(cols[5].get_text(strip=True) or 0),
                        days=d, start_time=(times[i].split("-")[0] if i < len(times) and "-" in times[i] else None),
                        end_time=(times[i].split("-")[1] if i < len(times) and "-" in times[i] else None),
                        room=rooms[i] if i < len(rooms) else None,
                        lecturers=json.dumps(lecs) if lecs else None
                    ))
            return data
        except Exception as e:
            print(f"Error scraping {code}: {e}")
            return []


def insert_courses(conn, courses):
    with conn.cursor() as cur:
        cur.executemany("""
            INSERT INTO courses (course_code,title,lab_section,lec_section,room,credit,
                                 days,start_time,end_time,lecturers)
            VALUES (%(course_code)s,%(title)s,%(lab_section)s,%(lec_section)s,%(room)s,
                    %(credit)s,%(days)s,%(start_time)s,%(end_time)s,%(lecturers)s)
            ON CONFLICT DO NOTHING;
        """, courses)
    conn.commit()


# ==============================================================
# Job Runner
# ==============================================================

def scrape_all_courses(start, end, workers, batch_size, db, job_id):
    scraper = CourseScraperPool()
    codes   = [f"{i:06d}" for i in range(start, end+1)]
    with ThreadPoolExecutor(max_workers=workers) as ex, tqdm(total=len(codes), desc=f"Job {job_id}") as bar:
        batch = []
        for i, results in enumerate(ex.map(scraper.scrape_course, codes), 1):
            batch.extend(results)
            bar.update(1)
            if i % batch_size == 0 and batch:
                insert_courses(db, batch); batch.clear()
                update_job(db, job_id, "running", i)
        if batch: insert_courses(db, batch)
    scraper.cleanup()


def get_job(conn):
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        cur.execute("SELECT * FROM scrape_course_jobs WHERE status='pending' ORDER BY id LIMIT 1")
        return cur.fetchone()


def update_job(conn, job_id, status, progress=None):
    with conn.cursor() as cur:
        cur.execute(
            "UPDATE scrape_course_jobs SET status=%s, progress=COALESCE(%s,progress), updated_at=now() WHERE id=%s",
            (status, progress, job_id)
        )
    conn.commit()


def worker_loop(db_conf):
    conn = psycopg2.connect(**db_conf)
    while True:
        job = get_job(conn)
        if job:
            print(f"Starting job {job['id']}...")
            update_job(conn, job['id'], "running", 0)
            try:
                scrape_all_courses(job['start_code'], job['end_code'], job['workers'], 100, conn, job['id'])
                update_job(conn, job['id'], "completed", job['total'])
                print(f"Job {job['id']} completed")
            except Exception as e:
                print("Job failed:", e); update_job(conn, job['id'], "failed")
        time.sleep(10)


if __name__ == "__main__":
    print("Course scraper started...")
    worker_loop(dict(dbname="blueprint", user="melkey", password="password1234", host="localhost", port=5432))

