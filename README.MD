# G Bridge
## Libary for Google Smart Home Actions
**-- work in progress --** 

Google smart home lets users control your connected devices through the Google Home app and the Google Assistant, which is available on more than 1 billion devices, like smart speakers, phones, cars, TVs, headphones, watches, and more. To connect your device with the Google Assistant, you need to build a smart home Action. Assistant handles how users trigger your Action (in multiple languages) and provides you with relevant metadata through Home Graph (such as a specific device based on user’s room); all you need to do is respond to the requests in your service.

This libary helps you to create a backend service for your Google Smart Home Action. 

This libary is ported from [actionssdk-smart-home-nodejs](https://github.com/actions-on-google/actionssdk-smart-home-nodejs)

# Philosopy
-- TODO: Explain how this libary is structured.. e.g duck-typing and interfaces usage--

### Useful links
* [https://developers.google.com/actions/smarthome/](https://developers.google.com/actions/smarthome/)
* [https://github.com/actions-on-google/actionssdk-smart-home-nodejs](https://github.com/actions-on-google/actionssdk-smart-home-nodejs)
* [https://console.actions.google.com/](https://console.actions.google.com/)

### Instructions 
##### (these may not be up to date, google changed their process recently)
1. Create an action.json file
```
    {
      "actions": [{
        "name": "actions.devices",
        "deviceControl": {
        },
        "fulfillment": {
          "conversationName": "automation"
        }
      }],
      "conversations": {
        "automation" :
        {
          "name": "automation",
          "url": "https://<YOUR URL>/smarthome"
        }
      }
    }
```

2. Create a project on [https://console.actions.google.com/](https://console.actions.google.com/)

3. Click **Use Actions SDK**

4. Use the [gActions](https://developers.google.com/actions/tools/gactions-cli) CLI to run the command given with the 'action.json' file as your Action Package.

5. Click **Okay**.

6. Click **ADD** under **App information**.

7. Give your App some information like an invocation name, some description, and some policy and contact info.

8. Click **Save**.

9. Click **Add** under **Account Linking**.

10. Select **Authorization Code** for Grant Type.

11. Under **Client Information**, enter the client ID and secret from earlier.

12. The Authorization URL is the hosted URL of your app with '/oauth' as the path, e.g. https://&lt;YOUR URL&gt;/oauth

13. The Token URL is the hosted URL of your app with '/token' as the path, e.g. https://&lt;YOUR URL&gt;/token

14. Click **TEST DRAFT**

15. Open the **Google Home App** on your phone and navigate to the **Home Control** section

15. Click the 3 dots and go to **Manage Accounts**

16. Select your project from the **Add New** section


### GoDoc 
--- TODO: Link to GODOC ----

### Example
Please see `/cmd/gbridge/` for examples on how to use this libary. 


# Notes
* If you push an update to actions.json, you might have to unlink your account in google home and re-link it for testing to work from your phone

* Place this service behind caddy or some other reverse proxy that does LetsEncrypt SSL, or do your own SSL encryption. Google needs your service to be reachable via HTTPs and a valid SSL certificate.  
