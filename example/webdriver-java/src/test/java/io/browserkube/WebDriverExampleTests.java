package io.browserkube;


import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;
import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.remote.LocalFileDetector;
import org.openqa.selenium.remote.RemoteWebDriver;
import org.openqa.selenium.remote.RemoteWebElement;

import java.io.IOException;
import java.net.URL;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.Properties;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class WebDriverExampleTests {
    Properties props;
    @BeforeAll
    void init() throws IOException {
        System.setProperty("webdriver.chrome.driver", "drivers/chromedriver-mac-64bit");
        this.props = readProps();
    }
    /**
     * Rigorous Test :-)
     */
    @Test
    public void shouldUploadAFile() throws IOException {
//        WebDriver driver = new ChromeDriver();
        ChromeOptions chromeOptions = new ChromeOptions();

        WebDriver driver = new RemoteWebDriver(new URL(props.getProperty("WD_URL")),
                chromeOptions);
        driver.get("https://ps.uci.edu/~franklin/doc/file_upload.html");

        String filePath = Thread.currentThread().getContextClassLoader().getResource("testupload.txt").getPath();
        //Locating upload filebutton
        WebElement uploadElement = driver.findElement(By.xpath("//input[@name='userfile']"));
        if (uploadElement instanceof RemoteWebElement) {
            ((RemoteWebElement) uploadElement).setFileDetector(new LocalFileDetector());
        }
        uploadElement.sendKeys(filePath);
        driver.quit();
    }

    @Test
    public void basicTest() throws Exception {
        Map<String, Object> browserkubeOptions = new HashMap<>();
        browserkubeOptions.put("enableVNC", true);
        browserkubeOptions.put("tenant", "test");
        browserkubeOptions.put("name", "java test session name");
        ChromeOptions chromeOptions = new ChromeOptions();
        chromeOptions.setCapability("browserName", "chrome");
        chromeOptions.setCapability("browserkube:options", browserkubeOptions);


        WebDriver driver = new RemoteWebDriver(new URL(props.getProperty("WD_URL")), chromeOptions);

        driver.get("https://go.dev/play/?simple=1");

        WebElement codeElement = driver.findElement(By.cssSelector("#code"));
        Assertions.assertNotNull(codeElement);

        codeElement.clear();
        codeElement.sendKeys("""
                    package main
                    import "fmt"
                    func main() {
                        fmt.Println("Hello WebDriver!")
                    }
                """
        );
        WebElement runButton = driver.findElement(By.cssSelector("#run"));
        Assertions.assertNotNull(runButton);
        runButton.click();

        WebElement outputDiv = driver.findElement(By.cssSelector("pre.Playground-output"));
        Assertions.assertNotNull(outputDiv);

        while (true) {
            String output = outputDiv.getText();
            if (!Objects.equals(output, "Waiting for remote server...")) {
                break;
            }
            Thread.sleep(250);
        }
        WebElement outputPre = driver.findElement(By.cssSelector("span.stdout"));
        Assertions.assertNotNull(outputDiv);

        String output = outputPre.getText();
        if (output.isEmpty()){
            Assertions.fail("output is empty");
        }
        System.out.println(output);

        driver.quit();
    }

    public static  Properties readProps() throws IOException {
            Properties props = new Properties();
            props.load(Thread.currentThread().getContextClassLoader().getResourceAsStream(".env"));
            return props;
    }
}
