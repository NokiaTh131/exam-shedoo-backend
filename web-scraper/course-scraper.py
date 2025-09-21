import argparse
from playwright.sync_api import sync_playwright
from bs4 import BeautifulSoup
import pandas as pd
from tqdm import tqdm
from concurrent.futures import ThreadPoolExecutor
import threading
import psycopg2
import json

class CourseScraperPool:
    def __init__(self, max_workers=4):
        self.max_workers = max_workers
        self.local_data = threading.local()
    
    def get_browser_page(self):
        if not hasattr(self.local_data, 'page'):
            self.local_data.playwright = sync_playwright().start()
            self.local_data.browser = self.local_data.playwright.chromium.launch(headless=True)
            self.local_data.page = self.local_data.browser.new_page()
            self.local_data.page.goto("https://www1.reg.cmu.ac.th/registrationoffice/searchcourse.php")
        return self.local_data.page
    
    def cleanup_browser(self):
        if hasattr(self.local_data, 'browser'):
            self.local_data.browser.close()
            self.local_data.playwright.stop()
    
    def scrape_course(self, course_code: str):
        try:
            page = self.get_browser_page()
            page.fill("#fcourse", "")
            page.fill("#fcourse", course_code)
            page.click("#button2")
            page.wait_for_load_state("networkidle", timeout=10000)
            
            html = page.content()
            soup = BeautifulSoup(html, "html.parser")
            table = soup.select_one(".tblCourse")
            
            data = []
            if table:
                rows = table.find_all("tr")
                for row in rows:
                    cols = row.find_all("td")
                    if len(cols) < 12:
                        continue
                    
                    section = cols[3].get_text(strip=True)
                    if not section.isdigit():
                        continue

                    day_cell = cols[7]
                    time_cell = cols[8]
                    room_cell = cols[9]
                    lecturer_cell = cols[10]

                    rooms = []
                    for element in room_cell.find_all(["div", "span"]):
                        text = element.get_text(strip=True)
                        if text and text != "-":
                            rooms.append(text)

                    days = []
                    for element in day_cell.find_all(["div", "span"]):
                        text = element.get_text(strip=True)
                        if text and text != "-":
                            days.append(text)

                    times = []
                    for element in time_cell.find_all(["div", "span"]):
                        text = element.get_text(strip=True)
                        if text and text != "-":
                            times.append(text)

                    lecturers = []
                    for element in lecturer_cell.stripped_strings:
                        text = element.strip()
                        lecturers.append(text)

                    for i in range(len(days)):
                        record = {
                            "course_code": course_code,
                            "title": cols[2].get_text(strip=True),
                            "lec_section": cols[3].get_text(strip=True) or None,
                            "lab_section": cols[4].get_text(strip=True) or None,
                            "credit": float(cols[5].get_text(strip=True) or 0),
                            "days": days[i] if i < len(days) else None,
                            "start_time": (times[i].split("-")[0] if "-" in times[i] else None),
                            "end_time": (times[i].split("-")[1] if "-" in times[i] else None),
                            "room": rooms[i] if i < len(rooms) else None,
                            "lecturers": lecturers or None,
                        }
                        data.append(record)
            return data
            
        except Exception as e:
            print(f"Error scraping {course_code}: {e}")
            return []


def worker_task(args):
    scraper, course_code = args
    return scraper.scrape_course(course_code)



def insert_courses(conn, courses):
    for c in courses:
        if isinstance(c.get("lecturers"), list):
            c["lecturers"] = json.dumps(c["lecturers"])
    
    with conn.cursor() as cur:
        cur.executemany(
            """
            INSERT INTO courses (
                course_code, title, lab_section, lec_section, room, credit,
                days, start_time, end_time, lecturers
            )
            VALUES (
                %(course_code)s, %(title)s, %(lab_section)s, %(lec_section)s, %(room)s,
                %(credit)s, %(days)s, %(start_time)s, %(end_time)s, %(lecturers)s
            )
            ON CONFLICT DO NOTHING;
            """,
            courses
        )
    conn.commit()

def scrape_all_courses_threaded(start=0, end=999999, max_workers=4, batch_size=100, db_config=None):
    all_data = []
    course_codes = [f"{i:06d}" for i in range(start, end + 1)]
    
    scraper = CourseScraperPool(max_workers)
    print(f"Starting scrape with {max_workers} workers for {len(course_codes)} courses...")

    conn = psycopg2.connect(**db_config)

    try:
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            tasks = [(scraper, code) for code in course_codes]
            
            with tqdm(total=len(course_codes), desc="Scraping courses") as pbar:
                batch_data = []
                
                for i, result in enumerate(executor.map(worker_task, tasks)):
                    if result:
                        batch_data.extend(result)
                    
                    pbar.update(1)
                    
                    if (i + 1) % batch_size == 0 and batch_data:
                        insert_courses(conn, batch_data)
                        all_data.extend(batch_data)
                        batch_data = []
                        print(f"\nInserted {len(all_data)} records so far...")
                
                if batch_data:
                    insert_courses(conn, batch_data)
                    all_data.extend(batch_data)
    
    finally:
        try:
            scraper.cleanup_browser()
        except:
            pass
        conn.close()
    
    print(f"\nCompleted! Scraped {len(all_data)} course records.")
    return all_data


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Scrape CMU courses by course code range")
    parser.add_argument("--start", type=int, required=True, help="Start course code (integer)")
    parser.add_argument("--end", type=int, required=True, help="End course code (integer)")
    parser.add_argument("--workers", type=int, default=4, help="Number of concurrent workers")
    parser.add_argument("--dbname", type=str, required=True, help="Postgres DB name")
    parser.add_argument("--user", type=str, required=True, help="Postgres user")
    parser.add_argument("--password", type=str, required=True, help="Postgres password")
    parser.add_argument("--host", type=str, default="localhost", help="Postgres host")
    parser.add_argument("--port", type=int, default=5432, help="Postgres port")
    args = parser.parse_args()
    
    db_config = {
        "dbname": args.dbname,
        "user": args.user,
        "password": args.password,
        "host": args.host,
        "port": args.port,
    }
    
    scrape_all_courses_threaded(
        start=args.start,
        end=args.end,
        max_workers=args.workers,
        db_config=db_config
    )

