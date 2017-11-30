2fa is a two-factor authentication manager in command line for macOS,
utilizing the macOS keychain to protect your TOTP secrets.

2fa was forked from [rsc/2fa](https://github.com/rsc/2fa) and rewrote.

Usage:

    go get -u github.com/hSATAC/2fa

    2fa add [account]
    2fa add --screenshot
    2fa list
    2fa show [account]

`2fa add [account]` adds a new key to the 2fa keychain with the given name. It
prints a prompt to standard error and reads a two-factor key from standard
input. Two-factor keys are short case-insensitive strings of letters A-Z and
digits 2-7.

`2fa add --screenshot` adds a new key by taking a screenshot of the qrcode.
See it in action here: [(gif)](http://ash.cat/5fMaVA).

The new key generates time-based (TOTP) authentication codes.

`2fa list` lists the names of all the keys in the keychain.

`2fa show [account]` prints a two-factor authentication code from the key with the
given account.

The default time-based authentication codes are derived from a hash of the
key and the current time, so it is important that the system clock have at
least one-minute accuracy.

## Example

During GitHub 2FA setup, at the “Scan this barcode with your app” step,
click the “enter this text code instead” link. A window pops up showing
“your two-factor secret,” a short string of letters and digits.

Add it to 2fa under the name github, typing the secret at the prompt:

    $ 2fa add github
    2fa key for github: JBSWY3DPEHPK3PXP
    $

Then whenever GitHub prompts for a 2FA code, run 2fa to obtain one:

    $ 2fa show github
          [14]  899393

The first 2 digits is the countdown of the TOTP.

Press any key to copy the code and exit.

## Requirement

Screenshot QRCode scanner requires [zbar](https://github.com/ZBar/ZBar) to run,
Please install zbar.

`$ brew install zbar`
