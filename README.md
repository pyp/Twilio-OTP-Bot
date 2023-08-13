<p align="center">
  <img src="https://media.discordapp.net/attachments/1047681400972791859/1140154407221207080/S8c6WAT.png" alt="Logo" width="30%">
</p>
<h1 align="center">Twillio OTP Bot</h1>
<p align="center">
  <em>$20/day who??? sincerely, @doozle/@pyp</em>
</p>

## "how"

- Requires: [Go](https://go.dev/).
- Create a [Twillio](https://www.twilio.com/en-us) free trial account.
- Optional/Recommended: Purchase Twillio Premium or else you're going to get <em>"Press any number on the keypad to run your code!"</em>
- Find your Twillio **auth token**, **account sid**, and the **phone number** that's Twillio gives you. (if they don't get one)
- Install [Ngrok](https://ngrok.com/download) and run `ngrok http 8080`. Copy the Ngrok URL that's given to you.
- Fill out lines `13-18` & `26-28` in **main.go**.
- Start: `go run .` Install all necessary modules.
- <em>For additional support, please use this [resource](https://www.google.com/). I am not helping anyone unless $$$.</em> 

## "what"

- Go Coding Language (nobody has made an otp bot in go on github)
- You can easily create your own scripts by adding to the dict in **internal/twillio/twillio.go**
- Features that most paid  OTP bots lack. (i.e. call-back, target-name, digits-amnt, customizable scripts). 
- Simple & easy to setup (if you have >=2iq)

## "useless"

<details>
<summary><strong>Changelog</strong></summary>
<br>

```diff
v0.0.1 ⋮ 08/12/2023
! Initial release
```
</details>
