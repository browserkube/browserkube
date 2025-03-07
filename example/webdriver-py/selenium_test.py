from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webdriver import WebDriver



def attach_to_session(command_executor, options, session_id):
    # Define the command for fetching the session
    original_execute = WebDriver.execute

    def new_command_executor(_, command, params=None):
        if command == "newSession":
            # If it is a "newSession" command, replace with the "session_id"
            return {'value': {'sessionId': session_id}}
        return original_execute(_, command, params)

    # Set the remote connection with the specified executor
    webdriver.remote.webdriver.WebDriver.execute = new_command_executor
    driver = webdriver.Remote(command_executor=command_executor, options=options)
    driver.session_id = session_id
    webdriver.remote.webdriver.WebDriver.execute = original_execute

    return driver

if __name__ == '__main__':
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

    session_id = driver.session_id

    driver2 = attach_to_session('http://kubernetes.docker.internal/browserkube/wd/hub',options, session_id)
    driver2.session_id = session_id
    driver2.get("https://www.python.org")
    print(driver.current_url)

    driver2.quit()
    print("quit")