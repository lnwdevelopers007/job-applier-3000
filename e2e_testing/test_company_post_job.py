"""Selenium test for company application accept/reject"""
from decouple import config
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait, Select
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
from construct_webdriver import get_driver
import time

URL = config('URL')
SEARCH_KEYWORD = config('TARGET_APPLICANT')
HEADLESS = config('HEADLESS', default='False') == 'True'
WAIT_TIME = config('WAIT_TIME', default=5, cast=int)
PROFILE = config('CHROME_COMPANY')
NUMBER = config('PROFILE_MUMBER_COMPANY', default="Profile 1")
JOBTITLE = config('JOBTITLE', default="Selenium Test Job")
JOBDESC = config('JOBDESCRIPTION', default="This is a test job posted by Selenium.")
JOBSUMMARY = config('JOBSUMMARY', default="Selenium Test Job Summary")
YEARS_EXPERIENCE = config('YEARSOFEXPERIENCE', default="2")
EDUCATIONLEVEL = config('EDUCATIONLEVEL', default="Bachelor's Degree")
LOCATION = config('LOCATION', default="Quarantine State")
POSTOPENDATE = config('POSTOPENDATE', default="12/12/2024")
POSTCLOSEDATE = config('POSTCLOSEDATE', default="01/31/2026")
MAXSALARY = config('MAXSALARY', default="100000")
MINSALARY = config('MINSALARY', default="50000")
EDITTITLE = config('EDITEDJOBTITLE', default="Updated Selenium Test Job")
DELETEREASON = config('DELETEREASON', default="Position filled")
SKILL = config('SKILL', default="JS")
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

    except:
        print("No login button found")

    createpost_nav = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//a[contains(., 'Post Job') or contains(., 'Post Job')]"))
    )
    createpost_nav.click()
    time.sleep(WAIT_TIME)

    print("Navigated to Post Job page")

    jobtitle_input = wait.until(
        EC.presence_of_element_located((By.XPATH, "//input[@placeholder='Enter job title...']"))
    )
    jobtitle_input.send_keys(JOBTITLE)

    location_input = driver.find_element(By.XPATH, "//input[@placeholder='Enter location...']")
    location_input.send_keys(LOCATION)
    max_salary_input = driver.find_element(By.XPATH, "//input[@placeholder='Enter maximum salary']")
    max_salary_input.send_keys(MAXSALARY)
    min_salary_input = driver.find_element(By.XPATH, "//input[@placeholder='Enter minimum salary']")
    min_salary_input.send_keys(MINSALARY)

    print("Filled first part of the job post form")

    continue_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Continue')]"))
    )
    continue_btn.click()
    time.sleep(WAIT_TIME)

    job_desc_input = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//div[@role='textbox']"))
    )
    job_desc_input.click()
    job_desc_input.send_keys(JOBDESC)

    job_summary_input = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//textarea"))
    )
    job_summary_input.send_keys(JOBSUMMARY)

    time.sleep(WAIT_TIME)    
    print("Filled second part of the job post form")

    continue_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Continue')]"))
    )
    continue_btn.click()

    years_select = Select(wait.until(
        EC.element_to_be_clickable((By.ID, "yearsOfExperience"))
    ))
    years_select.select_by_visible_text(YEARS_EXPERIENCE)

    education_select = Select(wait.until(
        EC.element_to_be_clickable((By.ID, "educationLevel"))
    ))
    education_select.select_by_visible_text(EDUCATIONLEVEL)

    skill_input = wait.until(
        EC.element_to_be_clickable((By.ID, "requiredSkills"))
    )
    skill_input.send_keys(SKILL)
    skill_input.send_keys(Keys.RETURN)
    time.sleep(WAIT_TIME)
    print("Filled third part of the job post form")

    continue_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Continue')]"))
    )
    continue_btn.click()

    posting_open_input = wait.until(
        EC.element_to_be_clickable((By.ID, "postingOpenDate"))
    )
    posting_open_input.clear()
    posting_open_input.send_keys(POSTOPENDATE)

    posting_close_input = wait.until(
        EC.element_to_be_clickable((By.ID, "postingCloseDate"))
    )
    posting_close_input.clear()
    posting_close_input.send_keys(POSTCLOSEDATE)

    print("Filled posting open/close dates")
    time.sleep(WAIT_TIME)

    submit_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Publish')]"))
    )
    submit_btn.click()
    print("Job posted successfully")
    time.sleep(WAIT_TIME)

    dashboard_nav = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//a[contains(., 'Dashboard') or contains(., 'Dashboard')]"))
    )
    dashboard_nav.click()
    time.sleep(WAIT_TIME)
    print("Navigated to Dashboard to verify the posted job")
    try:
        # Find the table row that contains the job title
        job_row = wait.until(
            EC.presence_of_element_located(
                (By.XPATH, f"//tr[.//*[contains(text(), '{JOBTITLE}')]]")
            )
        )

        # Inside that row, find the Edit button
        edit_btn = job_row.find_element(
            By.XPATH, ".//button[contains(., 'Edit')]"
        )
        edit_btn.click()
        print(f"Clicked Edit button for job '{JOBTITLE}'")
        time.sleep(WAIT_TIME)

        jobtitle_input = wait.until(
            EC.presence_of_element_located((By.XPATH, "//input[@placeholder='Enter job title...']"))
        )
        jobtitle_input.clear()
        jobtitle_input.send_keys(EDITTITLE)
        print(f"Changed job title to '{EDITTITLE}'")
        time.sleep(WAIT_TIME)

        submit_btn = wait.until(
            EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Save Changes')]"))
        )
        submit_btn.click()
        print(f"Job '{JOBTITLE}' edited successfully to '{EDITTITLE}'")
        time.sleep(WAIT_TIME)
    except:
        print(f"Edit button for job '{JOBTITLE}' NOT found")
    job_row = wait.until(
        EC.presence_of_element_located(
            (By.XPATH, f"//tr[.//*[contains(text(), '{EDITTITLE}')]]")
        )
    )
    delete_btn = job_row.find_element(
        By.XPATH, ".//button[contains(., 'Delete')]"
    )
    delete_btn.click()
    print(f"Clicked Delete button for job '{EDITTITLE}'")
    delete_input = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//textarea"))
    )
    delete_input.send_keys(DELETEREASON)
    time.sleep(WAIT_TIME)
    confirm_delete_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Confirm Delete')]"))
    )
    confirm_delete_btn.click()
    time.sleep(WAIT_TIME)
    print(f"Job '{EDITTITLE}' deleted successfully")
except:
    print("Company job posting test failed")