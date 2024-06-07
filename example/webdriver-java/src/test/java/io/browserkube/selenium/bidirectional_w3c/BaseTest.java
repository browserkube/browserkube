package io.browserkube.selenium.bidirectional_w3c;

import io.browserkube.WebDriverExampleTests;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.remote.AbstractDriverOptions;
import org.openqa.selenium.remote.RemoteWebDriver;

import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.Map;
import java.util.Properties;

public class BaseTest {
    protected static Properties props;
    protected WebDriver driver;
    private ChromeOptions chromeOptions = new ChromeOptions();

    @BeforeAll
    static void init() throws IOException {
        System.setProperty("webdriver.chrome.driver", "drivers/chromedriver-mac-64bit");
        props = WebDriverExampleTests.readProps();
    }

    @BeforeEach
    void setup() throws MalformedURLException {
        driver = (RemoteWebDriver) RemoteWebDriver.builder()
                .address(new URL(props.getProperty("WD_URL")))
                .addAlternative(getOptionsFor(chromeOptions))
                .build();
    }

    @AfterEach
    void quit() {
        if (driver != null) {
            driver.quit();
        }
    }

    private AbstractDriverOptions getOptionsFor(AbstractDriverOptions options) {
        options.setCapability("webSocketUrl", true);
        options.setCapability("browserName", "chrome");
        options.setCapability("browserkube:options", Map.of("enableVNC", true, "tenant", "test"));
        options.setCapability("webSocketUrl", true);
        return options;
    }
}