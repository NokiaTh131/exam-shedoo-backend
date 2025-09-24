import time, requests, psycopg2, traceback, os
from psycopg2.extras import RealDictCursor
from bs4 import BeautifulSoup
from dotenv import load_dotenv

EXAM_TYPES = ["MIDTERM", "FINAL"]

# ------------------------------
# Helper functions
# ------------------------------
def parse_time_range(time_str):
    if "-" in time_str:
        start, end = time_str.split("-")
        return start.strip(), end.strip()
    return time_str, None

def get_pending_job(conn):
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        cur.execute("SELECT * FROM scrape_exam_jobs WHERE status='pending' ORDER BY id LIMIT 1")
        return cur.fetchone()

def update_job(conn, job_id, status):
    with conn.cursor() as cur:
        cur.execute(
            "UPDATE scrape_exam_jobs SET status=%s, updated_at=NOW() WHERE id=%s",
            (status, job_id)
        )
    conn.commit()

# ------------------------------
# Scraper
# ------------------------------
def scrape_and_insert(term, conn):
    cursor = conn.cursor()
    for exam_type in EXAM_TYPES:
        url = f"https://www1.reg.cmu.ac.th/registrationoffice/timetable_exam.php?type={exam_type}&term={term}"
        r = requests.get(url)
        r.encoding = "utf-8"
        if r.status_code != 200:
            print(f"Failed to fetch {exam_type}, status {r.status_code}")
            continue

        soup = BeautifulSoup(r.text, "html.parser")
        table = soup.find("table")
        if not table:
            print(f"No exam table found for {exam_type}")
            continue

        time_slots = [th.get_text(strip=True) for th in table.find("thead").find_all("tr")[1].find_all("th")]

        for tr in table.find("tbody").find_all("tr"):
            cols = tr.find_all("td")
            if not cols: 
                continue
            date = cols[0].get_text(strip=True)

            for i, td in enumerate(cols[1:]):
                start_time, end_time = parse_time_range(time_slots[i])
                courses = [c.strip() for c in td.get_text().split(",") if c.strip()]
                for course in courses:
                    if course.upper() == "REGULAR EXAM": 
                        continue

                    # Fetch course_id first
                    cursor.execute("SELECT id FROM courses WHERE course_code=%s LIMIT 1", (course,))
                    result = cursor.fetchone()
                    if not result:
                        print(f"Skipping {course} because it's not in courses table")
                        continue
                    course_id = result[0]

                    if exam_type == "MIDTERM":
                        cursor.execute("""
                            INSERT INTO course_exams(course_id, course_code, lec_section, lab_section, midterm_exam_date, midterm_exam_start_time, midterm_exam_end_time)
                            VALUES (%s, %s, '000', '000', %s, %s, %s)
                            ON CONFLICT(course_code, lec_section, lab_section)
                            DO UPDATE SET
                                midterm_exam_date=EXCLUDED.midterm_exam_date,
                                midterm_exam_start_time=EXCLUDED.midterm_exam_start_time,
                                midterm_exam_end_time=EXCLUDED.midterm_exam_end_time
                        """, (course_id, course, date, start_time, end_time))
                    else:
                        cursor.execute("""
                            INSERT INTO course_exams(course_id, course_code, lec_section, lab_section, final_exam_date, final_exam_start_time, final_exam_end_time)
                            VALUES (%s, %s, '000', '000', %s, %s, %s)
                            ON CONFLICT(course_code, lec_section, lab_section)
                            DO UPDATE SET
                                final_exam_date=EXCLUDED.final_exam_date,
                                final_exam_start_time=EXCLUDED.final_exam_start_time,
                                final_exam_end_time=EXCLUDED.final_exam_end_time
                        """, (course_id, course, date, start_time, end_time))
    conn.commit()
    cursor.close()

# ------------------------------
# Worker loop
# ------------------------------
def worker_loop(db_config: dict):
    print("Exam scraper worker started...")
    conn = psycopg2.connect(**db_config)
    while True:
        job = get_pending_job(conn)
        if job:
            job_id, term = job["id"], job["term"]
            print(f"Starting job {job_id} (term={term})...")
            update_job(conn, job_id, "running")
            try:
                scrape_and_insert(term, conn)
                update_job(conn, job_id, "completed")
                print(f"Job {job_id} completed")
            except Exception as e:
                print(f"Job {job_id} failed:", e)
                update_job(conn, job_id, "failed")
        time.sleep(5)

# ------------------------------
# Entry point
# ------------------------------
if __name__ == "__main__":
    print("Course scraper started...")

    # Load environment variables from .env file
    load_dotenv()

    db_config = dict(
        dbname=os.getenv("POSTGRES_DATABASE"),
        user=os.getenv("POSTGRES_USERNAME"),
        password=os.getenv("POSTGRES_PASSWORD"),
        host=os.getenv("POSTGRES_HOST"),
        port=int(os.getenv("POSTGRES_PORT", 5432))  # default to 5432 if not set
    )

    worker_loop(db_config)
