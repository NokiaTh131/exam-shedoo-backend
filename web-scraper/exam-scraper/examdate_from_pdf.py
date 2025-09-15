import re
import pdfplumber
import psycopg2
import argparse

# --- Regex patterns ---
thai_date_pattern = re.compile(r'(\d{1,2})\s*([ก-ฮ\.]+)\s*(\d{2,4})')
time_pattern = re.compile(r'(\d{2}[:.]?\d{2}[-–—]\d{2}[:.]?\d{2})')
course_pattern = re.compile(r'\b(\d{6})\b')

# Map Thai months to English abbreviations for consistency
thai_to_eng_month = {
    'ม.ค.':'JAN','ก.พ.':'FEB','มี.ค.':'MAR','เม.ย.':'APR','พ.ค.':'MAY','มิ.ย.':'JUN',
    'ก.ค.':'JUL','ส.ค.':'AUG','ก.ย.':'SEP','ต.ค.':'OCT','พ.ย.':'NOV','ธ.ค.':'DEC'
}

def parse_thai_date_raw(day, month_abbr, year):
    """Return date in format like 'OCT 20'"""
    month_eng = thai_to_eng_month.get(month_abbr.strip())
    if not month_eng:
        return None
    return f"{month_eng} {int(day):02d}"

def parse_pdf_and_insert(pdf_path, db_config, exam_type="MIDTERM"):
    results = []

    # --- Extract data from PDF ---
    with pdfplumber.open(pdf_path) as pdf:
        current_date = None
        current_time = None
        for page in pdf.pages:
            text = page.extract_text() or ""
            for line in text.splitlines():
                line = line.replace("\xa0", " ").replace("\u202f", " ").strip()

                dmatch = thai_date_pattern.search(line)
                if dmatch:
                    day, month_abbr, year = dmatch.groups()
                    current_date = parse_thai_date_raw(day, month_abbr, year)

                tmatch = time_pattern.search(line)
                if tmatch:
                    current_time = tmatch.group(1)

                cmatch = course_pattern.search(line)
                if cmatch:
                    course_code = cmatch.group(1)
                    results.append((current_date, current_time, course_code))

    # --- Insert into PostgreSQL ---
    conn = psycopg2.connect(
        host=db_config["host"],
        port=db_config["port"],
        dbname=db_config["dbname"],
        user=db_config["user"],
        password=db_config["password"]
    )
    cursor = conn.cursor()

    for date, time_str, course_code in results:
        if not date or not time_str:
            continue

        start_time, end_time = time_str.split("-")
        start_time = start_time.replace(":", "").strip()
        end_time = end_time.replace(":", "").strip()

        if exam_type.upper() == "MIDTERM":
            cursor.execute("""
                INSERT INTO course_exams (
                    course_code,
                    midterm_exam_date,
                    midterm_exam_start_time,
                    midterm_exam_end_time
                ) VALUES (%s, %s, %s, %s)
                ON CONFLICT (course_code, section) 
                DO UPDATE SET
                    midterm_exam_date = EXCLUDED.midterm_exam_date,
                    midterm_exam_start_time = EXCLUDED.midterm_exam_start_time,
                    midterm_exam_end_time = EXCLUDED.midterm_exam_end_time;
            """, (course_code, date, start_time, end_time))
        else:  # FINAL
            cursor.execute("""
                INSERT INTO course_exams (
                    course_code,
                    final_exam_date,
                    final_exam_start_time,
                    final_exam_end_time
                ) VALUES (%s, %s, %s, %s)
                ON CONFLICT (course_code, section) 
                DO UPDATE SET
                    final_exam_date = EXCLUDED.final_exam_date,
                    final_exam_start_time = EXCLUDED.final_exam_start_time,
                    final_exam_end_time = EXCLUDED.final_exam_end_time;
            """, (course_code, date, start_time, end_time))

    conn.commit()
    cursor.close()
    conn.close()
    print(f"Data from {pdf_path} inserted into PostgreSQL successfully.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Parse Thai PDF exam schedule and insert into PostgreSQL")
    parser.add_argument("--pdf", required=True, help="Path to PDF file")
    parser.add_argument("--exam_type", required=True, choices=["MIDTERM", "FINAL"], help="Type of exam")
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

    parse_pdf_and_insert(args.pdf, db_config, args.exam_type)

