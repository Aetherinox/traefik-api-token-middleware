<div align="center">
<h6>Traefik v3 middleware which allows for you to protect certain aspects of your site with an API key.</h6>
<h2>♾️ Traefik API Key Middleware ♾️</h1>

<br />

<p>

This Traefik middleware allows you to secure certain routes behind a request header API key. Users who have not successfully authenticated will be greeted with a **403 Forbidden Error**.

</p>

<br />

<img src="https://github.com/user-attachments/assets/9e4b02a7-1a77-4175-b7c5-149b272d29f7" height="230">

<br />
<br />

</div>

<div align="center">

<!-- prettier-ignore-start -->
[![Version][github-version-img]][github-version-uri]
[![Downloads][github-downloads-img]][github-downloads-uri]
[![Size][github-size-img]][github-size-img]
[![Last Commit][github-commit-img]][github-commit-img]
[![Contributors][contribs-all-img]](#contributors-)

[![Built with Material for MkDocs](https://img.shields.io/badge/Powered_by_Material_for_MkDocs-526CFE?style=for-the-badge&logo=MaterialForMkDocs&logoColor=white)](https://aetherinox.github.io/traefik-api-token-middleware/)
<!-- prettier-ignore-end -->

</div>

<br />

---

<br />

- [Configuration](#configuration)
  - [Static File](#static-file)
    - [File (YAML)](#file-yaml)
    - [File (TOML)](#file-toml)
    - [CLI](#cli)
  - [Dynamic File](#dynamic-file)
    - [File (YAML)](#file-yaml-1)
    - [File (TOML)](#file-toml-1)
    - [Kubernetes Custom Resource Definition](#kubernetes-custom-resource-definition)
- [Parameters](#parameters)
  - [authenticationErrorMsg](#authenticationerrormsg)
- [Usage](#usage)
- [Browser Plugins](#browser-plugins)
  - [Firefox](#firefox)
    - [Extension: Header Editor](#extension-header-editor)
    - [Extension: Modify Header Value](#extension-modify-header-value)
- [Verifying Modified Headers](#verifying-modified-headers)
- [Contributors ✨](#contributors-)

<br />

---

<br />

## Configuration
The following provides examples for usage scenarios.

<br />

### Static File
If you are utilizing a Traefik **Static File**, review the following examples:

<br />

#### File (YAML)

```yaml
## Static configuration
experimental:
  plugins:
    traefik-api-token-middleware:
      moduleName: "github.com/Aetherinox/traefik-api-token-middleware"
      version: "v0.1.0"
```

<br />

#### File (TOML)

```toml
## Static configuration
[experimental.plugins.traefik-api-token-middleware]
  moduleName = "github.com/Aetherinox/traefik-api-token-middleware"
  version = "v0.1.0"
```

<br />

#### CLI

```bash
## Static configuration
--experimental.plugins.traefik-api-token-middleware.modulename=github.com/Aetherinox/traefik-api-token-middleware
--experimental.plugins.traefik-api-token-middleware.version=v0.1.0
```

<br />

### Dynamic File
If you are utilizing a Traefik **Dynamic File**, review the following examples:

<br />

#### File (YAML)

```yaml
# Dynamic configuration
http:
  middlewares:
    api-token:
      plugin:
        traefik-api-token-middleware:
          authenticationHeader: true
          authenticationHeaderName: X-API-TOKEN
          bearerHeader: true
          bearerHeaderName: Authorization
          removeHeadersOnSuccess: true
          tokens:
            - your-api-token
```

<br />

#### File (TOML)

```toml
# Dynamic configuration
[http]
  [http.middlewares]
    [http.middlewares.api-token]
      [http.middlewares.api-token.plugin]
        [http.middlewares.api-token.plugin.traefik-api-token-middleware]
          authenticationHeader = true
          authenticationHeaderName = "X-API-TOKEN"
          bearerHeader = true
          bearerHeaderName = "Authorization"
          removeHeadersOnSuccess = true
          tokens = ["your-api-token"]
```

<br />

#### Kubernetes Custom Resource Definition

```yaml
# Dynamic configuration
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: api-token
spec:
  plugin:
    traefik-api-token-middleware:
      authenticationHeader: true
      authenticationHeaderName: X-API-TOKEN
      bearerHeader: true
      bearerHeaderName: Authorization
      removeHeadersOnSuccess: true
      tokens:
        - your-api-token
```

<br />

---

<br />

## Parameters
This plugin accepts the following parameters:

| Parameter | Description | Default | Type | Required |
| --- | --- | --- | --- | --- |
| `authenticationHeader` | Pass token using Authentication Header | true | bool | | 
| `authenticationHeaderName` | Authentication header name | 'X-API-TOKEN' | string | | 
| `authenticationErrorMsg` | Error message to display on unsuccessful authentication | 'Access Denied' | string | |
| `bearerHeader` | Pass token using Authentication Header Bearer Key | true | bool | |
| `bearerHeaderName` | Authentication bearer header name | 'Authorization' | string | |
| `removeHeadersOnSuccess` | If `true`; remove header on successful authentication | true | bool | |
| `noPublicTokenOnError` | Don't display name of token in unsuccessful error message | false | bool | |
| `timestampUnix` | Display datetime in Unix timestamp instead of UnixDate | false | bool | |

<br />

### authenticationErrorMsg
This setting changes how error messages are displayed to a user who doesn't provide a correct token. If `enabled`, it will keep the name of your token private.

<br />

---

<br />

## Usage
Implementing an API key route is simple.

<br />

---

<br />

## Browser Plugins
If you do not want to specify an API key using conventional means (such as by using `curl`), you can utilize a front-end browser extension. This allows you to supply a modified request header with your specific API key which will grant you access to your desired location.

<br />

### Firefox
If you are using Firefox, install the plugin(s) below. ( Pick **One** ):

- [Extension: Header Editor](#extension-header-editor)
- [Extension: Modify Header Value](#modify-header-value)

<br />

#### Extension: Header Editor  
With this extension, you can modify the request header and response header, cancel a request and redirect a request.
**Please do not write regular expressions that begin with `(. *), (. *?), (. +)`, such regular expressions may cause problems with Firefox**  

<br />

Once you install the browser extension above, open the **settings**.

<br />

<p align="center"><img style="width: 50%;text-align: center;" src="https://github.com/user-attachments/assets/309b841d-1c74-447a-82d5-a4a508811a8a"></p>

<br />

Create a new rule
- Name: `API Token`
- Rule Type: `Modify request header`
- Match Type: `Domain`
    - Match Rules: `subdomain.domain.com`
    - Exclude Rule: `none`
- Execute Type: `normal`
    - Header Name: `x-api-token`
    - Header Value: `your-api-token`

<br />

<p align="center"><img style="width: 80%;text-align: center;" src="https://github.com/user-attachments/assets/313a6cfd-6416-4fc4-acd5-3cfe2a44c502"></p>

<br />

Once you have your modified header added to the browser extension, verify it by reading the section [Verifying Modified Headers](#verifying-modified-headers).

<br />
<br />

#### Extension: Modify Header Value
[Modify Header Value](https://prod.outgoing.prod.webservices.mozgcp.net/v1/6eb6158dedc247c12e9010ccc61bb16738f844da3ac6765a5f43f7eb5ebace6f/https%3A//mybrowseraddon.com/modify-header-value.html) can add, modify or remove an HTTP-request-header for all requests on a desired website or URL. This Firefox add-on is very useful if you are an App developer, website designer, or if you want to test a particular header for a request on a website.

<br />

Once you install the browser extension above, open the **settings**.

<br />

<p align="center"><img style="width: 50%;text-align: center;" src="https://github.com/user-attachments/assets/50d4d8fb-aee1-4aa8-8f67-db645dd2909d"></p>

<br />

You need to add a new rule which injects your modified header into the specified domain.

<br />

<p align="center"><img style="width: 100%;text-align: center;" src="https://github.com/user-attachments/assets/df389731-0bcd-4b39-a143-9ba49cef22f4"></p>

<br />

Once you have your modified header added to the browser extension, verify it by reading the section [Verifying Modified Headers](#verifying-modified-headers).

<br />

---

<br />

## Verifying Modified Headers
This section explains how you can verify that your modified request header is being received by your server.

If you now go to the subdomain and open the browser developer tools, and on the **Inspector** tab, scroll to the very bottom.

<br />



<br />

---

<br />

## Contributors ✨
We are always looking for contributors. If you feel that you can provide something useful to Gistr, then we'd love to review your suggestion. Before submitting your contribution, please review the following resources:

- [Pull Request Procedure](.github/PULL_REQUEST_TEMPLATE.md)
- [Contributor Policy](CONTRIBUTING.md)

<br />

Want to help but can't write code?
- Review [active questions by our community](https://github.com/Aetherinox/traefik-api-token-middleware/labels/help%20wanted) and answer the ones you know.

<br />

The following people have helped get this project going:

<br />

<div align="center">

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![Contributors][contribs-all-img]](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top"><a href="https://gitlab.com/Aetherinox"><img src="https://avatars.githubusercontent.com/u/118329232?v=4?s=40" width="80px;" alt="Aetherinox"/><br /><sub><b>Aetherinox</b></sub></a><br /><a href="https://github.com/Aetherinox/traefik-api-token-middleware/commits?author=Aetherinox" title="Code">💻</a> <a href="#projectManagement-Aetherinox" title="Project Management">📆</a> <a href="#fundingFinding-Aetherinox" title="Funding Finding">🔍</a></td>
    </tr>
  </tbody>
</table>
</div>
<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

<br />
<br />

<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<!-- BADGE > GENERAL -->
  [general-npmjs-uri]: https://npmjs.com
  [general-nodejs-uri]: https://nodejs.org
  [general-npmtrends-uri]: http://npmtrends.com/traefik-api-token-middleware

<!-- BADGE > VERSION > GITHUB -->
  [github-version-img]: https://img.shields.io/github/v/tag/Aetherinox/traefik-api-token-middleware?logo=GitHub&label=Version&color=ba5225
  [github-version-uri]: https://github.com/Aetherinox/traefik-api-token-middleware/releases

<!-- BADGE > VERSION > NPMJS -->
  [npm-version-img]: https://img.shields.io/npm/v/traefik-api-token-middleware?logo=npm&label=Version&color=ba5225
  [npm-version-uri]: https://npmjs.com/package/traefik-api-token-middleware

<!-- BADGE > VERSION > PYPI -->
  [pypi-version-img]: https://img.shields.io/pypi/v/traefik-api-token-middleware-plugin
  [pypi-version-uri]: https://pypi.org/project/traefik-api-token-middleware-plugin/

<!-- BADGE > LICENSE > MIT -->
  [license-mit-img]: https://img.shields.io/badge/MIT-FFF?logo=creativecommons&logoColor=FFFFFF&label=License&color=9d29a0
  [license-mit-uri]: https://github.com/Aetherinox/traefik-api-token-middleware/blob/main/LICENSE

<!-- BADGE > GITHUB > DOWNLOAD COUNT -->
  [github-downloads-img]: https://img.shields.io/github/downloads/Aetherinox/traefik-api-token-middleware/total?logo=github&logoColor=FFFFFF&label=Downloads&color=376892
  [github-downloads-uri]: https://github.com/Aetherinox/traefik-api-token-middleware/releases

<!-- BADGE > NPMJS > DOWNLOAD COUNT -->
  [npmjs-downloads-img]: https://img.shields.io/npm/dw/%40aetherinox%2Ftraefik-api-token-middleware?logo=npm&&label=Downloads&color=376892
  [npmjs-downloads-uri]: https://npmjs.com/package/traefik-api-token-middleware

<!-- BADGE > GITHUB > DOWNLOAD SIZE -->
  [github-size-img]: https://img.shields.io/github/repo-size/Aetherinox/traefik-api-token-middleware?logo=github&label=Size&color=59702a
  [github-size-uri]: https://github.com/Aetherinox/traefik-api-token-middleware/releases

<!-- BADGE > NPMJS > DOWNLOAD SIZE -->
  [npmjs-size-img]: https://img.shields.io/npm/unpacked-size/traefik-api-token-middleware/latest?logo=npm&label=Size&color=59702a
  [npmjs-size-uri]: https://npmjs.com/package/traefik-api-token-middleware

<!-- BADGE > CODECOV > COVERAGE -->
  [codecov-coverage-img]: https://img.shields.io/codecov/c/github/Aetherinox/traefik-api-token-middleware?token=MPAVASGIOG&logo=codecov&logoColor=FFFFFF&label=Coverage&color=354b9e
  [codecov-coverage-uri]: https://codecov.io/github/Aetherinox/traefik-api-token-middleware

<!-- BADGE > ALL CONTRIBUTORS -->
  [contribs-all-img]: https://img.shields.io/github/all-contributors/Aetherinox/traefik-api-token-middleware?logo=contributorcovenant&color=de1f6f&label=contributors
  [contribs-all-uri]: https://github.com/all-contributors/all-contributors

<!-- BADGE > GITHUB > BUILD > NPM -->
  [github-build-img]: https://img.shields.io/github/actions/workflow/status/Aetherinox/traefik-api-token-middleware/npm-release.yml?logo=github&logoColor=FFFFFF&label=Build&color=%23278b30
  [github-build-uri]: https://github.com/Aetherinox/traefik-api-token-middleware/actions/workflows/npm-release.yml

<!-- BADGE > GITHUB > BUILD > Pypi -->
  [github-build-pypi-img]: https://img.shields.io/github/actions/workflow/status/Aetherinox/traefik-api-token-middleware/release-pypi.yml?logo=github&logoColor=FFFFFF&label=Build&color=%23278b30
  [github-build-pypi-uri]: https://github.com/Aetherinox/traefik-api-token-middleware/actions/workflows/pypi-release.yml

<!-- BADGE > GITHUB > TESTS -->
  [github-tests-img]: https://img.shields.io/github/actions/workflow/status/Aetherinox/traefik-api-token-middleware/npm-tests.yml?logo=github&label=Tests&color=2c6488
  [github-tests-uri]: https://github.com/Aetherinox/traefik-api-token-middleware/actions/workflows/npm-tests.yml

<!-- BADGE > GITHUB > COMMIT -->
  [github-commit-img]: https://img.shields.io/github/last-commit/Aetherinox/traefik-api-token-middleware?logo=conventionalcommits&logoColor=FFFFFF&label=Last%20Commit&color=313131
  [github-commit-uri]: https://github.com/Aetherinox/traefik-api-token-middleware/commits/main/

<!-- prettier-ignore-end -->
<!-- markdownlint-restore -->