"""Construct a webdriver"""
from selenium import webdriver
from selenium.webdriver.chrome.options import Options

def get_driver(profile, number, headless=True):
    """Construct driver"""
    options = Options()
    if headless:
        options.add_argument("--headless")
        options.add_argument("--no-sandbox")
    else:
        options.add_argument("--start-maximized")
        options.add_argument(f"--user-data-dir={profile}")
        options.add_argument(f"--profile-directory={number}")
    options.add_argument("--disable-gpu")
    driver = webdriver.Chrome(options=options)
    return driver