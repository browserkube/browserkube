const http = require("http")
const dotenv = require("dotenv")

dotenv.config()

const capabilities = JSON.stringify({
	"desiredCapabilities": {
		"browserName": "chrome",
		"browserkube:options": {
			"enableVideo": true
		},
		"browserVersion": "116.0",
		"enableVNC": true,
		"saveVideoEndpoint": "file:///home/seluser/videos",
		"labels": {
			"manual": "true"
		},
		"sessionTimeout": "60m",
		"name": "Manual session"
	},
	"capabilities": {
		"alwaysMatch": {
			"browserName": "chrome",
			"browserVersion": "116.0",
			"browserkube:options": {
				"enableVideo": true
			},
			"selenoid:options": {
				"enableVNC": true,
				"sessionTimeout": "60m",
				"saveVideoEndpoint": "file:///home/seluser/videos",
				"labels": {
					"manual": "true"
				},
				"screenResolution": "1920x1080x24"
			},
			"goog:chromeOptions": {
				"args": [
					"start-maximized"
				]
			}
		},
		"firstMatch": [
			{}
		]
	}
});

function startBrowser() {
    return new Promise((resolve, reject) => {
		const host = process.env.BROWSERKUBE_URL;
		const port = process.env.BROWSERKUBE_PORT;
        var browserkubeEndpointOpts = {
            host: host,
            port: port,
            path: '/browserkube/wd/hub/session?timeout=60000',
            method: 'POST',
            headers: {
                Accept: 'application/json,test/plain,*/*',
                Connection: 'keep-alive'
            }
        };
        const req = http.request(browserkubeEndpointOpts, function(res) {
            let error;
            console.log('Browser Creation Status: ' + res.statusCode);
            if (res.statusCode !== 200) {
                error = new Error('Request Failed');
            }
            if (error) {
                console.error(error.message);
                res.resume();
                return reject(error);
            }
            let rawData = '';
            res.setEncoding('utf8')
            res.on('data', function (chunk) {
                console.log('Session Create Body: ' + chunk);
                rawData += chunk;
            })
            res.on('end', () => {
                try {
                    resolve(JSON.parse(rawData));
                } catch (e) {
                    console.error(e.message);
                }
            });
        });
        
        req.write(capabilities);
        req.end();
    });
}

function stopBrowser(sessionID) {
	return new Promise((resolve, reject) => {
		const host = process.env.BROWSERKUBE_URL;
		const port = process.env.BROWSERKUBE_PORT;
		var deleteBrowserOpts = {
			host: host,
            port: port,
            path: '/browserkube/wd/hub/session/' + sessionID + '?timeout=50000',
            method: 'DELETE',
		}
		const req = http.request(deleteBrowserOpts, function(res) {
			let error;
			console.log('Browser Deletion Status: ' + res.statusCode);
			if (res.statusCode !== 200) {
				error = new Error('Request Failed')
			}
			if (error) {
				console.error(error.message);
				res.resume();
				return reject(error);
			}
		}).end();
	})
}

module.exports = {startBrowser, stopBrowser};