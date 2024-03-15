# gwurl - "Google distributing Windows installer"'s URL

[![CI - Nix Status](https://github.com/kachick/gwurl/actions/workflows/ci-nix.yml/badge.svg?branch=main)](https://github.com/kachick/gwurl/actions/workflows/ci-nix.yml?query=branch%3Amain+)

## Usage

```console
> gwurl 'https://dl.google.com/tag/s/appguid%3D%7BDDCCD2A9-025E-4142-BCEB-F467B88CF830%7D%26iid%3D%7BBF6725E7-4A15-87B7-24D5-0ADABEC753EE%7D%26lang%3Dja%26browser%3D4%26usagestats%3D0%26appname%3DGoogle%2520%25E6%2597%25A5%25E6%259C%25AC%25E8%25AA%259E%25E5%2585%25A5%25E5%258A%259B%26needsadmin%3Dtrue%26ap%3Dexternal-stable-universal/japanese-ime/GoogleJapaneseInputSetup.exe'

{appguid:{DDCCD2A9-025E-4142-BCEB-F467B88CF830} ap:external-stable-universal appname:Google 日本語入力 needsadmin:true filename:GoogleJapaneseInputSetup.exe}
http://edgedl.me.gvt1.com/edgedl/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi
https://edgedl.me.gvt1.com/edgedl/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi
http://dl.google.com/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi
https://dl.google.com/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi
http://www.google.com/dl/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi
https://www.google.com/dl/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi
```

## Tested targets

- Google Chrome
- Google Japanese IME(Also known as "Google 日本語入力")

## Installation

Nix

```bash
nix run github:kachick/gwurl -- 'the_tagged_url_here'
```

Prebuilt binaries

```pwsh
# <UPDATE ME>
```

## Motivation

As my understanding.

If we are downloading Google product as a Windows installer, it returns long and encoded URL as `https://dl.google.com/tag/s/appguid%3D%7~iid%3D%`.\
The iid is the "Instance ID", so downloading from this URL twice, the installer will not work. Even if you got the binary...!\
As far as I know, they does not open the list of these actual URLs.

I didn't know their motivation, but this is a blocker of getting installer permalinks. Especially in [winget-pkgs](https://github.com/microsoft/winget-pkgs/pull/144281#discussion_r1524718296).

## Resources

Using their API without any token auth and referenced these OSS resources.

- https://github.com/google/omaha/blob/c97bc28a89fa0d9863a7e9089c37357eb854e8de/doc/TaggedMetainstallers.md?plain=1#L11-L21
- https://github.com/google/omaha/blob/c97bc28a89fa0d9863a7e9089c37357eb854e8de/doc/ServerProtocolV2.md?plain=1#L21
- https://github.com/SpecterShell/Dumplings/blob/26a3b4f359a0fb524dc9f30dc5ae924c0f723f8a/Tasks/Google.Chrome/Script.ps1

## Plans

Current features are enough to me, but it would be good...

- Output as JSON: For automated GitHub Actions to create PRs for winget-pkgs
- Returns supported minimum Windows Versions: Also helps to create winget-pkgs
