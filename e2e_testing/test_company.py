"""Selenium test for company application accept/reject"""
from decouple import config
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
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
ACCEPT_OR_REJECT = config("ACTION")

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

    applicants_nav = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//a[contains(., 'Applicants') or contains(., 'Applicants')]"))
    )
    applicants_nav.click()
    time.sleep(WAIT_TIME)

    # Search for a job
    search_bar = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//input[contains(@placeholder, 'Search') or @type='text']"))
    )
    search_bar.clear()
    search_bar.send_keys(SEARCH_KEYWORD)
    time.sleep(WAIT_TIME)

    pending_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(., 'Pending')]"))
    )
    pending_btn.click()
    print("Pending button clicked")
    time.sleep(WAIT_TIME)

    first_applicant = wait.until(
        EC.element_to_be_clickable(
            (By.XPATH, "/html/body/div/div[1]/div/div/main/div/div[2]/div[1]/div[2]/button[1]/div[3]/h3")
        )
    )
    first_applicant.click()
    print("First applicant (Pending) clicked")
    time.sleep(WAIT_TIME)

    action_xpath = f"//button[(contains(@class, 'bg-green-600') or contains(@class, 'bg-red-600')) and contains(normalize-space(.), '{ACCEPT_OR_REJECT}')]"
    action_btn = wait.until(
        EC.element_to_be_clickable((By.XPATH, action_xpath))
    )
    action_btn.click()
    print(f"{ACCEPT_OR_REJECT} button clicked")
    time.sleep(WAIT_TIME)

    print("Accept/Reject completed!")

except:
    print("Company accept/reject failed")