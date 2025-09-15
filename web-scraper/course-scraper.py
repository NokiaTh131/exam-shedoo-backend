from playwright.sync_api import sync_playwright
from bs4 import BeautifulSoup
import pandas as pd
from tqdm import tqdm
from concurrent.futures import ThreadPoolExecutor
import threading

class CourseScraperPool:
    def __init__(self, max_workers=4):
        self.max_workers = max_workers
        self.local_data = threading.local()
    
    def get_browser_page(self):
        """Get or create browser and page for current thread"""
        if not hasattr(self.local_data, 'page'):
            self.local_data.playwright = sync_playwright().start()
            self.local_data.browser = self.local_data.playwright.chromium.launch(headless=True)
            self.local_data.page = self.local_data.browser.new_page()
            self.local_data.page.goto("https://www1.reg.cmu.ac.th/registrationoffice/searchcourse.php")
        return self.local_data.page
    
    def cleanup_browser(self):
        """Clean up browser resources for current thread"""
        if hasattr(self.local_data, 'browser'):
            self.local_data.browser.close()
            self.local_data.playwright.stop()
    
    def scrape_course(self, course_code: str):
        """
        Scrape a single course by code. Returns a list of dicts (rows) or empty list if no data.
        """
        try:
            page = self.get_browser_page()
            
            # Clear and fill the course code
            page.fill("#fcourse", "")
            page.fill("#fcourse", course_code)
            page.click("#button2")
            
            # Wait for response with timeout
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
                    
                    record = {
                        "course_id": course_code,
                        "course_title": cols[2].get_text(strip=True),
                        "lec_section": cols[3].get_text(strip=True),
                        "lab_section": cols[4].get_text(strip=True),
                        "lec_credit": cols[5].get_text(strip=True),
                        "lab_credit": cols[6].get_text(strip=True),
                        "day": cols[7].get_text(strip=True),
                        "time": cols[8].get_text(strip=True).rstrip("-"),
                        "room": cols[9].get_text(strip=True),
                        "lecturer": cols[10].get_text(strip=True),
                    }
                    data.append(record)
            
            return data
            
        except Exception as e:
            print(f"Error scraping {course_code}: {e}")
            return []

def worker_task(args):
    """Worker function for thread pool"""
    scraper, course_code = args
    return scraper.scrape_course(course_code)

def scrape_all_courses_threaded(start=0, end=999999, save_csv="all_courses.csv", max_workers=4, batch_size=100):
    """
    Scrape courses using ThreadPoolExecutor for improved performance
    
    Args:
        start: Starting course number
        end: Ending course number  
        save_csv: Output CSV filename
        max_workers: Number of concurrent threads
        batch_size: Save progress every N courses
    """
    all_data = []
    course_codes = [f"{i:06d}" for i in range(start, end + 1)]
    
    # Create scraper instance
    scraper = CourseScraperPool(max_workers)
    
    print(f"Starting scrape with {max_workers} workers for {len(course_codes)} courses...")
    
    try:
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            # Prepare arguments for workers
            tasks = [(scraper, code) for code in course_codes]
            
            # Process with progress bar
            with tqdm(total=len(course_codes), desc="Scraping courses") as pbar:
                batch_data = []
                
                for i, result in enumerate(executor.map(worker_task, tasks)):
                    if result:
                        batch_data.extend(result)
                    
                    pbar.update(1)
                    
                    # Save progress periodically
                    if (i + 1) % batch_size == 0:
                        all_data.extend(batch_data)
                        df_temp = pd.DataFrame(all_data)
                        df_temp.to_csv(f"temp_{save_csv}", index=False)
                        batch_data = []
                        print(f"\nSaved progress: {len(all_data)} records so far...")
                
                # Add remaining data
                all_data.extend(batch_data)
    
    finally:
        # Clean up resources
        try:
            scraper.cleanup_browser()
        except:
            pass
    
    # Save final results
    df = pd.DataFrame(all_data)
    df.to_csv(save_csv, index=False)
    
    # Clean up temp file
    try:
        import os
        os.remove(f"temp_{save_csv}")
    except:
        pass
    
    print(f"\nCompleted! Scraped {len(all_data)} course records.")
    return df

if __name__ == "__main__":
    df = scrape_all_courses_threaded()
    print("Scraped all courses")
