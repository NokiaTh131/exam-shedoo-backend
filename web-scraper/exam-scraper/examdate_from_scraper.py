import argparse
import requests
from bs4 import BeautifulSoup
import psycopg2

EXAM_TYPES = ["MIDTERM", "FINAL"]

def parse_time_range(time_str):
    """Convert '0800-1100' -> ('0800', '1100')"""
    if "-" in time_str:
        start, end = time_str.split("-")
        return start.strip(), end.strip()
    return time_str, None

def scrape_and_insert(term_arg, db_config):
    conn = psycopg2.connect(
        host=db_config["host"],
        port=db_config["port"],
        dbname=db_config["dbname"],
        user=db_config["user"],
        password=db_config["password"]
    )
    cursor = conn.cursor()

    for exam_type in EXAM_TYPES:
        url = f"https://www1.reg.cmu.ac.th/registrationoffice/timetable_exam.php?type={exam_type}&term={term_arg}"
        response = requests.get(url)
        response.encoding = "utf-8"

        if response.status_code != 200:
            print(f"Failed to fetch {exam_type}, status code: {response.status_code}")
            continue

        soup = BeautifulSoup(response.text, "html.parser")
        table = soup.find("table")
        if not table:
            print(f"No exam timetable table found for {exam_type}.")
            continue

        # Extract time slots from thead
        thead_rows = table.find("thead").find_all("tr")
        time_slots = [th.get_text(strip=True) for th in thead_rows[1].find_all("th")]

        # Extract data from tbody
        for tr in table.find("tbody").find_all("tr"):
            cols = tr.find_all("td")
            if not cols:
                continue

            date = cols[0].get_text(strip=True)

            for i, td in enumerate(cols[1:]):
                time_range = time_slots[i]
                start_time, end_time = parse_time_range(time_range)

                courses = [c.strip() for c in td.get_text().split(",") if c.strip()]
                for course in courses:
                    if course.upper() == "REGULAR EXAM":
                        continue

                    if exam_type == "MIDTERM":
                        cursor.execute("""
                            INSERT INTO course_exams (
                                course_code,
                                midterm_exam_date,
                                midterm_exam_start_time,
                                midterm_exam_end_time
                            ) VALUES (%s, %s, %s, %s)
                            ON CONFLICT (course_code) 
                            DO UPDATE SET
                                midterm_exam_date = EXCLUDED.midterm_exam_date,
                                midterm_exam_start_time = EXCLUDED.midterm_exam_start_time,
                                midterm_exam_end_time = EXCLUDED.midterm_exam_end_time;
                        """, (course, date, start_time, end_time))
                    else:  # FINAL
                        cursor.execute("""
                            INSERT INTO course_exams (
                                course_code,
                                final_exam_date,
                                final_exam_start_time,
                                final_exam_end_time
                            ) VALUES (%s, %s, %s, %s)
                            ON CONFLICT (course_code) 
                            DO UPDATE SET
                                final_exam_date = EXCLUDED.final_exam_date,
                                final_exam_start_time = EXCLUDED.final_exam_start_time,
                                final_exam_end_time = EXCLUDED.final_exam_end_time;
                        """, (course, date, start_time, end_time))

    conn.commit()
    cursor.close()
    conn.close()
    print("Data inserted into PostgreSQL successfully.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Scrape CMU exams and insert into PostgreSQL")
    parser.add_argument("--term", required=True, help="Term ID, e.g. 168")
    parser.add_argument("--host", required=True, help="PostgreSQL host")
    parser.add_argument("--port", required=True, type=int, help="PostgreSQL port")
    parser.add_argument("--dbname", required=True, help="PostgreSQL database name")
    parser.add_argument("--user", required=True, help="PostgreSQL username")
    parser.add_argument("--password", required=True, help="PostgreSQL password")
    args = parser.parse_args()

    db_config = {
        "host": args.host,
        "port": args.port,
        "dbname": args.dbname,
        "user": args.user,
        "password": args.password
    }

    scrape_and_insert(args.term, db_config)

