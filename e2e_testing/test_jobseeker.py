"""Selenium test for jobseeker application submission"""
from decouple import config
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
from construct_webdriver import get_driver
import time

URL = config('URL')
SEARCH_KEYWORD = config('TARGET_JOB_SEARCH')
HEADLESS = config('HEADLESS', default='False') == 'True'
WAIT_TIME = config('WAIT_TIME', default=5, cast=int)
PROFILE = config('CHROME_JOBSEEKER')
NUMBER = config('PROFILE_MUMBER_JOBSEEKER', default="Profile 1")

driver = get_driver(PROFILE, NUMBER, HEADLESS)
actions = ActionChains(driver)
wait = WebDriverWait(driver, 20)

try:
    # Open the URL
    driver.get(URL)
    print(f"Entering site: {URL}")

    time.sleep(WAIT_TIME)
    try:
        # Press Login button
        login_btn = wait.until(
            EC.element_to_be_clickable((By.XPATH, "//a[contains(text(), 'Log in') or contains(., 'Log in')]"))
        )
        login_btn.click()
        time.sleep(WAIT_TIME)

        # Login via Google
        google_btn = wait.until(
            EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Google') or contains(@aria-label, 'Google')]"))
        )
        google_btn.click()

        # Switch the window to Google login
        driver.switch_to.window(driver.window_handles[-1])
        time.sleep(WAIT_TIME)

        # Wait for redirect to job page
        wait.until(EC.url_contains("jobs"))
        print("Login successful")
        time.sleep(WAIT_TIME)

    except:
        print("No login button found")

    job_nav = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//a[contains(., 'Jobs') or contains(., 'Jobs')]"))
    )
    job_nav.click()
    time.sleep(WAIT_TIME)

    # Search for a job
    search_bar = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//input[contains(@placeholder, 'Search') or @type='text']"))
    )
    search_bar.clear()
    search_bar.send_keys(SEARCH_KEYWORD)
    time.sleep(WAIT_TIME)

    # Click on first job
    first_job_title = wait.until(
        EC.element_to_be_clickable((By.XPATH, "/html/body/div/div[1]/div/div/main/div/main/div[2]/section[1]/div/div/div[1]/div/div[1]/div/div[1]/h3"))
    )
    first_job_title.click()
    time.sleep(WAIT_TIME)

    # Apply for the first job in the list
    apply_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Apply') or contains(., 'Apply Now')]"))
    )
    apply_btn.click()
    time.sleep(WAIT_TIME)

    # Submit the application
    submit_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Submit Application')]"))
    )
    submit_btn.click()
    time.sleep(WAIT_TIME*2)

    print("Application submitted")

    applications_nav = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//a[contains(., 'Applications') or contains(., 'My Applications')]"))
    )
    applications_nav.click()
    time.sleep(WAIT_TIME)

    all_spans = wait.until(
    EC.presence_of_all_elements_located((By.TAG_NAME, "span"))
    )

    found = False
    for span in all_spans:
        if SEARCH_KEYWORD.lower() in span.text.lower():
            found = True
            print(f"Found application containing keyword: {span.text}")
            break

    if not found:
        print(f"No span contains the keyword: {SEARCH_KEYWORD}")

except:
    print("Job Application failed")