import re
import pdfplumber

PDF_PATH = "midterm_168.pdf"

thai_date_pattern = re.compile(r'(\d{1,2})\s*([ก-ฮ\.]+)\s*(\d{2,4})')
time_pattern = re.compile(r'(\d{2}[:.]?\d{2}[-–—]\d{2}[:.]?\d{2})')
course_pattern = re.compile(r'\b(\d{6})\b')

thai_months = {
    'ม.ค.':'01','ก.พ.':'02','มี.ค.':'03','เม.ย.':'04','พ.ค.':'05','มิ.ย.':'06',
    'ก.ค.':'07','ส.ค.':'08','ก.ย.':'09','ต.ค.':'10','พ.ย.':'11','ธ.ค.':'12'
}

def parse_thai_date(day, month_abbr, year):
    try:
        month = thai_months.get(month_abbr.strip())
        if not month:
            return None
        year = int(year)
        if year < 100:  # Two-digit year
            year += 2500
        ce_year = year - 543
        return f"{int(day):02d}-{month}-{ce_year:04d}"
    except:
        return None

results = []


with pdfplumber.open(PDF_PATH) as pdf:
    current_date = None
    current_time = None  # Track the last seen time
    for page in pdf.pages:
        text = page.extract_text() or ""
        for line in text.splitlines():
            line = line.replace("\xa0", " ").replace("\u202f", " ").strip()

            dmatch = thai_date_pattern.search(line)
            if dmatch:
                day, month_abbr, year = dmatch.groups()
                current_date = parse_thai_date(day, month_abbr, year)

            tmatch = time_pattern.search(line)
            if tmatch:
                current_time = tmatch.group(1)

            cmatch = course_pattern.search(line)
            if cmatch:
                course_code = cmatch.group(1)
                results.append((current_date, current_time, course_code))

for date, time_str, course_code in results:
    print(f"{date} {time_str} {course_code}")

