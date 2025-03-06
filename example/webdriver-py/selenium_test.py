from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys

options = webdriver.FirefoxOptions()
options.set_capability("browserkube:options", {
    "enableVNC": True,
})
# options.add_argument("--no-sandbox")
# options.add_argument(f"--user-data-dir=/home/user")

# Create a new instance of the Chrome driver
driver = webdriver.Remote(
    options=options,
    command_executor='http://kubernetes.docker.internal/browserkube/wd/hub')

# Open the Python website
driver.get("https://www.python.org")

# Print the page title
print(driver.title)

# Find the search bar using its name attribute
search_bar = driver.find_element(By.NAME, "q")
search_bar.clear()
search_bar.send_keys("getting started with python")
search_bar.send_keys(Keys.RETURN)

# Print the current URL
print(driver.current_url)

# Close the browser window
driver.close()
