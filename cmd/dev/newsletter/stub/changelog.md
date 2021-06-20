
## Breaking Changes


We listened to your feedback and have improved the naming of the SDK method `initializeSelfServiceRecoveryForNativeApps` to better match what it does: `initializeSelfServiceRecoveryWithoutBrowser`. As in the previous release you may still use the old SDK if you do not want to deal with the SDK breaking changes for now.
We listened to your feedback and have improved the naming of the SDK method `initializeSelfServiceVerificationForNativeApps` to better match what it does: `initializeSelfServiceVerificationWithoutBrowser`. As in the previous release you may still use the old SDK if you do not want to deal with the SDK breaking changes for now.
We listened to your feedback and have improved the naming of the SDK method `initializeSelfServiceSettingsForNativeApps` to better match what it does: `initializeSelfServiceSettingsWithoutBrowser`. As in the previous release you may still use the old SDK if you do not want to deal with the SDK breaking changes for now.
We listened to your feedback and have improved the naming of the SDK method `initializeSelfServiceregistrationForNativeApps` to better match what it does: `initializeSelfServiceregistrationWithoutBrowser`. As in the previous release you may still use the old SDK if you do not want to deal with the SDK breaking changes for now.
We listened to your feedback and have improved the naming of the SDK method `initializeSelfServiceLoginForNativeApps` to better match what it does: `initializeSelfServiceLoginWithoutBrowser`. As in the previous release you may still use the old SDK if you do not want to deal with the SDK breaking changes for now.

**Bug Fixes:**

* Add json detection to setting error subbranches ([fb83dcb](https://github.com/ory/kratos/commit/fb83dcb8ae7463079ddb33c04673cf4556f6058c))
* Cache migration status ([5be2f14](https://github.com/ory/kratos/commit/5be2f149cd79ddfbe8496eccf5d5aacb6a9a0b8e)), closes [#1337](https://github.com/ory/kratos/issues/1337)
* Change SMTP config validation from URI to a Regex pattern ([#1436](https://github.com/ory/kratos/issues/1436)) ([5ab1e8f](https://github.com/ory/kratos/commit/5ab1e8f17bcbc229fada2c584b2c1f576b819761)), closes [#1435](https://github.com/ory/kratos/issues/1435)
* Check filesystem before fallback to bundled templates ([#1401](https://github.com/ory/kratos/issues/1401)) ([22d999e](https://github.com/ory/kratos/commit/22d999e78eb4f67d2f3ba07e62fd28ffb3331d6d))
* Continue button for oidc registration step ([2aad5ac](https://github.com/ory/kratos/commit/2aad5ac8f7055f39f4f434d26fbca74cdbe75337)), closes [#1422](https://github.com/ory/kratos/issues/1422) [#1320](https://github.com/ory/kratos/issues/1320):

    When signing up with an OIDC provider and the traits model is missing some fields, the submit button shows all OIDC options. Instead, it should show just one option called "Continue".

* Deprecate sessionCookie ([#1428](https://github.com/ory/kratos/issues/1428)) ([eccad74](https://github.com/ory/kratos/commit/eccad741a1702181d4b207aad954a950906a808b)), closes [#1426](https://github.com/ory/kratos/issues/1426)
* Do not cache incomplete migrations ([#1434](https://github.com/ory/kratos/issues/1434)) ([154c26f](https://github.com/ory/kratos/commit/154c26f6da4bb7040deabdc352c90cdae42c69fe))
* Do not run network migrations when booting ([12bbab9](https://github.com/ory/kratos/commit/12bbab9d3cf788998cd4a9be50ac8c7a9d2232bd)), closes [#1399](https://github.com/ory/kratos/issues/1399)
* Improve identity list performance ([f76886f](https://github.com/ory/kratos/commit/f76886fe7436f71fbef00081888a2f8d0106ba98)), closes [#1412](https://github.com/ory/kratos/issues/1412)
* Incorrect openapi specification for verification submission  ([#1431](https://github.com/ory/kratos/issues/1431)) ([ecb0a01](https://github.com/ory/kratos/commit/ecb0a01f61441aa97751943b5e9ddcc28f783d91)), closes [#1368](https://github.com/ory/kratos/issues/1368)
* Mark ui node message as optional ([#1365](https://github.com/ory/kratos/issues/1365)) ([7b8d59f](https://github.com/ory/kratos/commit/7b8d59f48ed14a6d0672238645d8675d4bf7fd77)), closes [#1361](https://github.com/ory/kratos/issues/1361) [#1362](https://github.com/ory/kratos/issues/1362)
* Mark verified_at as omitempty ([77b258e](https://github.com/ory/kratos/commit/77b258e57a3d53fe437838a5e9c57805e9c970aa)):

    Closes https://github.com/ory/sdk/issues/46

* Panic if contextualizer is not set ([760035a](https://github.com/ory/kratos/commit/760035a6c5efa08561b93daff57ebb4655032b2a))
* Panic on error in issue session ([5fbd855](https://github.com/ory/kratos/commit/5fbd8557e1f907dd400bfcd26c187db16dc344ba)), closes [#1384](https://github.com/ory/kratos/issues/1384)
* Prometheus metrics fix ([#1299](https://github.com/ory/kratos/issues/1299)) ([ac5d00d](https://github.com/ory/kratos/commit/ac5d00d472a87ab51e7c6834e2cb59f107fc3b3b))
* Recovery email case sensitive ([#1357](https://github.com/ory/kratos/issues/1357)) ([bce14c4](https://github.com/ory/kratos/commit/bce14c487450bd668859f362b98704644fa4c72a)), closes [#1329](https://github.com/ory/kratos/issues/1329)
* Remove typing from node.attribute.value ([63a5e08](https://github.com/ory/kratos/commit/63a5e08afab76dafbfe13e6126e165af28492aad)):

    Closes https://github.com/ory/sdk/issues/75
    Closes https://github.com/ory/sdk/issues/74
    Closes https://github.com/ory/sdk/issues/72

* Rename client package for external consumption ([cba8b00](https://github.com/ory/kratos/commit/cba8b00c8b755cc0bdc7818bc9d7390ff3532ce1))
* Resolve driver issues ([47b1c8d](https://github.com/ory/kratos/commit/47b1c8dce57a023e89a2b178bc8a033496ef4ff2))
* Resolve network regression ([8f96b1f](https://github.com/ory/kratos/commit/8f96b1fe4d0846a3ad97a45bc972ece04109289d))
* Resolve network regressions ([8fc52c0](https://github.com/ory/kratos/commit/8fc52c034ed9978c2a04cc66bccc9b795c9bbefa))
* Testhelper regressions ([bf3b04f](https://github.com/ory/kratos/commit/bf3b04fd2c7f9162073cb584d6fb0d59e868ecbf))
* Use correct url in submitSelfServiceVerificationFlow ([ab8a600](https://github.com/ory/kratos/commit/ab8a600080ac0d6a6235806b74c5b9e3dc1c2d60))
* Use STARTTLS for smtps connections ([#1430](https://github.com/ory/kratos/issues/1430)) ([c21bb80](https://github.com/ory/kratos/commit/c21bb80a749df7b224a8ac3f15fa62523a78d805)), closes [#781](https://github.com/ory/kratos/issues/781)
* Version schema ([#1359](https://github.com/ory/kratos/issues/1359)) ([8c4bac7](https://github.com/ory/kratos/commit/8c4bac71674e45e440d916c6c947ed018a8ea29a)), closes [#1331](https://github.com/ory/kratos/issues/1331) [#1101](https://github.com/ory/kratos/issues/1101) [ory/hydra#2427](https://github.com/ory/hydra/issues/2427)

**Code Refactoring:**

* Corp package ([#1402](https://github.com/ory/kratos/issues/1402)) ([0202dc5](https://github.com/ory/kratos/commit/0202dc57aacc0d48e4c1ee4e68c91654451f63fa))
* Introduce DefaultContextualizer in corp package ([#1390](https://github.com/ory/kratos/issues/1390)) ([944d045](https://github.com/ory/kratos/commit/944d045aa7fc59eadfdd18951f0d4937b1ea79df)), closes [#1363](https://github.com/ory/kratos/issues/1363)
* Move cleansql to separate package ([7c203dc](https://github.com/ory/kratos/commit/7c203dc8219afe07f180143f832158615b51f60a))

**Documentation:**

* Add docs for registration SPA flow ([84458f1](https://github.com/ory/kratos/commit/84458f1a9dfe8be6a97bddd832fcc508b60b8498))
* Add go sdk examples ([e948fad](https://github.com/ory/kratos/commit/e948faddce3a1f52df964c701f6ba2a28f5dfe03))
* Add replit instructions ([8ab8607](https://github.com/ory/kratos/commit/8ab8607dee433f6e708ade296a6c26d0a87d0aae))
* Add tested and running go sdk examples ([3b56bb5](https://github.com/ory/kratos/commit/3b56bb5fd37d0e7d4479967aa0b5721a68a267f2))
* Fix typo in "Sign in/up with ID & assword" ([#1383](https://github.com/ory/kratos/issues/1383)) ([f39739d](https://github.com/ory/kratos/commit/f39739d94e97f20b94630b957371d11294dc8300))
* Mark login endpoints as experimental ([6faf0f6](https://github.com/ory/kratos/commit/6faf0f65bb05bbafdee6b1274a719695fd5b4173))
* Update docs for all flows ([d29ea69](https://github.com/ory/kratos/commit/d29ea69f6bb908b529502030942b1ced52227372))
* Update documentation for plaintext templates ([#1369](https://github.com/ory/kratos/issues/1369)) ([419784d](https://github.com/ory/kratos/commit/419784dd0d4ddc338830ed0d77a7d99f8f440777)), closes [#1351](https://github.com/ory/kratos/issues/1351)
* Update path ([f0384d9](https://github.com/ory/kratos/commit/f0384d9c11085230fd16290c524d22fac6002870))
* Update sdk use ([bcb8c06](https://github.com/ory/kratos/commit/bcb8c06ee324c639e548fc06315d9e952f470582))
* Use correct path ([#1333](https://github.com/ory/kratos/issues/1333)) ([e401135](https://github.com/ory/kratos/commit/e401135cf415d7e3e6a8ca463dd47e46fe399b33))

**Features:**

* Add GetContextualizer ([ac32717](https://github.com/ory/kratos/commit/ac3271742c9c2b968b08dd2b35a5d120c5befcd9))
* Add instana as possible tracing provider ([#1429](https://github.com/ory/kratos/issues/1429)) ([abe48a9](https://github.com/ory/kratos/commit/abe48a97ee75567979a70f00dd73ff698efcc75d)), closes [#1385](https://github.com/ory/kratos/issues/1385)
* Add vk and yandex providers to oidc providers and documentation ([#1339](https://github.com/ory/kratos/issues/1339)) ([22a3ef9](https://github.com/ory/kratos/commit/22a3ef98181eb5922cc0f1c016d42ce46732d0a2)), closes [#1234](https://github.com/ory/kratos/issues/1234)
* Improve contextualization in serve/daemon ([f83cd35](https://github.com/ory/kratos/commit/f83cd355422fb4b422f703406473bda914d8419c))
* Include Credentials Metadata in admin api ([#1274](https://github.com/ory/kratos/issues/1274)) ([c8b6219](https://github.com/ory/kratos/commit/c8b62190fca53db4e1b3a4ddb5253fbd2fd46002)), closes [#820](https://github.com/ory/kratos/issues/820)
* Include Credentials Metadata in admin api Missing changes in handler ([#1366](https://github.com/ory/kratos/issues/1366)) ([a71c220](https://github.com/ory/kratos/commit/a71c2208dedac45d32dab578e62a5e3105c8dee0))
* Natively support SPA for login flows ([6ff67af](https://github.com/ory/kratos/commit/6ff67afa8b0fc0a95cec44d3dda2cbc1987b51dd)), closes [#1138](https://github.com/ory/kratos/issues/1138) [#668](https://github.com/ory/kratos/issues/668):

    This patch adds the long-awaited capabilities for natively working with SPAs and AJAX requests. Previously, requests to the `/self-service/login/browser` endpoint would always end up in a redirect. Now, if the `Accept` header is set to `application/json`, the login flow will be returned as JSON instead. Accordingly, changes to the error and submission flow have been made to support `application/json` content types and SPA / AJAX requests.

* Natively support SPA for recovery flows ([5461244](https://github.com/ory/kratos/commit/5461244943286081e13c304a3b38413b8ee6fdf2)):

    This patch adds the long-awaited capabilities for natively working with SPAs and AJAX requests. Previously, requests to the `/self-service/recovery/browser` endpoint would always end up in a redirect. Now, if the `Accept` header is set to `application/json`, the registration flow will be returned as JSON instead. Accordingly, changes to the error and submission flow have been made to support `application/json` content types and SPA / AJAX requests.

* Natively support SPA for registration flows ([57d3c57](https://github.com/ory/kratos/commit/57d3c5786a88f0648e7fa57f181f060a057ec19f)), closes [#1138](https://github.com/ory/kratos/issues/1138) [#668](https://github.com/ory/kratos/issues/668):

    This patch adds the long-awaited capabilities for natively working with SPAs and AJAX requests. Previously, requests to the `/self-service/registration/browser` endpoint would always end up in a redirect. Now, if the `Accept` header is set to `application/json`, the registration flow will be returned as JSON instead. Accordingly, changes to the error and submission flow have been made to support `application/json` content types and SPA / AJAX requests.

* Natively support SPA for settings flows ([ea4395e](https://github.com/ory/kratos/commit/ea4395ed25d5668e4ce365336cd7a5e13e0ba1cc)):

    This patch adds the long-awaited capabilities for natively working with SPAs and AJAX requests. Previously, requests to the `/self-service/settings/browser` endpoint would always end up in a redirect. Now, if the `Accept` header is set to `application/json`, the registration flow will be returned as JSON instead. Accordingly, changes to the error and submission flow have been made to support `application/json` content types and SPA / AJAX requests.

* Natively support SPA for verification flows ([c151500](https://github.com/ory/kratos/commit/c1515009dcd1b5946a93733feedb01753de91c3d)):

    This patch adds the long-awaited capabilities for natively working with SPAs and AJAX requests. Previously, requests to the `/self-service/verification/browser` endpoint would always end up in a redirect. Now, if the `Accept` header is set to `application/json`, the registration flow will be returned as JSON instead. Accordingly, changes to the error and submission flow have been made to support `application/json` content types and SPA / AJAX requests.

* Sign in with Auth0 ([#1352](https://github.com/ory/kratos/issues/1352)) ([f618a53](https://github.com/ory/kratos/commit/f618a53fb971ad16121aa8728cfec54253bb3f44)), closes [#609](https://github.com/ory/kratos/issues/609)
* Support api in settings error ([23105db](https://github.com/ory/kratos/commit/23105dbb836d920b8766536b65de58932f53d6f6))
* Support reading session token from X-Session-Token HTTP header ([dcaefd9](https://github.com/ory/kratos/commit/dcaefd94a0b2cf819424f2e10b3bdae63b256726))
* Team id in slack oidc ([#1409](https://github.com/ory/kratos/issues/1409)) ([e4d021a](https://github.com/ory/kratos/commit/e4d021a037a6b44f8bd66372e9c260c640e87b9d)), closes [#1408](https://github.com/ory/kratos/issues/1408)
* Update openapi specs and regenerate ([cac507e](https://github.com/ory/kratos/commit/cac507eb5b1f39d003d72e57912dbbfe6f92deb1))
* **identities:** Add a state to identities ([#1312](https://github.com/ory/kratos/issues/1312)) ([d22954e](https://github.com/ory/kratos/commit/d22954e2fdb7b2dd5206651b6dd5cf96185a33ba)), closes [#598](https://github.com/ory/kratos/issues/598)

**Tests:**

* Add tests for cookie behavior of API and browser endpoints ([d1b1521](https://github.com/ory/kratos/commit/d1b15217867cfb92a615c793b26fad288f5e5742))
* Remove obsolete console.log ([3ecc869](https://github.com/ory/kratos/commit/3ecc869ebfef5c97334ae4334fb4af98ca9baf97))
* Resolve e2e regressions ([b0d3b82](https://github.com/ory/kratos/commit/b0d3b82f301942bebe3c0027c8b3160749f907af))
* Resolve migratest panic ([89d05ae](https://github.com/ory/kratos/commit/89d05ae0c376c4ea1f23708cccf95c9754a29c94))
* **e2e:** Greatly improve test performance ([#1421](https://github.com/ory/kratos/issues/1421)) ([2ffad9e](https://github.com/ory/kratos/commit/2ffad9ee751471451e2151719a2e70d5f89437b0)):

    Instead of running the individual profiles as separate Cypress instances, we now use one singular instance which updates the Ory Kratos configuration depending on the test context. This ensures that hot-reloading is properly working while also signficantly reducing the amount of time spent on booting up the service dependencies.


**Unclassified:**

* add CoC shield (#1439) ([826ed1a](https://github.com/ory/kratos/commit/826ed1a6deafdc2631a5c72f0bfacc91b06a3435)), closes [#1439](https://github.com/ory/kratos/issues/1439)
* u ([b03549b](https://github.com/ory/kratos/commit/b03549b6340ec0bf4f9d741ce145ca90bbc09968))
* Format ([5cc9fc3](https://github.com/ory/kratos/commit/5cc9fc3a6e91a96225d016d60c8da5cef647ac18))
* u ([318a31d](https://github.com/ory/kratos/commit/318a31d400b97653b4f377c67df4ae0afea189d9))
* Format ([e525805](https://github.com/ory/kratos/commit/e525805246431075d26c3f47596ae93f6580d8ee))
* Format ([4a692ac](https://github.com/ory/kratos/commit/4a692acc7db160068ed7d81461b173bc957e4736))
* Format ([169c0cd](https://github.com/ory/kratos/commit/169c0cd8d424babef69a52ddf65e2b75ded09a46))





## Breaking Changes


Unfortunately, some method signatures have changed in the SDKs. Below is a list of changed entries:

- Error `genericError` was renamed to `jsonError` and now includes more information and better typing for errors;
- The following functions have been renamed:
   - `initializeSelfServiceLoginViaAPIFlow` -> `initializeSelfServiceLoginForNativeApps`
   - `initializeSelfServiceLoginViaBrowserFlow` -> `initializeSelfServiceLoginForBrowsers`
   - `initializeSelfServiceRegistrationViaAPIFlow` -> `initializeSelfServiceRegistrationForNativeApps`
   - `initializeSelfServiceRegistrationViaBrowserFlow` -> `initializeSelfServiceRegistrationForBrowsers`
   - `initializeSelfServiceSettingsViaAPIFlow` -> `initializeSelfServiceSettingsForNativeApps`
   - `initializeSelfServiceSettingsViaBrowserFlow` -> `initializeSelfServiceSettingsForBrowsers`
   - `initializeSelfServiceRecoveryViaAPIFlow` -> `initializeSelfServiceRecoveryForNativeApps`
   - `initializeSelfServiceRecoveryViaBrowserFlow` -> `initializeSelfServiceRecoveryForBrowsers`
   - `initializeSelfServiceVerificationViaAPIFlow` -> `initializeSelfServiceVerificationForNativeApps`
   - `initializeSelfServiceVerificationViaBrowserFlow` -> `initializeSelfServiceVerificationForBrowsers`
- Some type names have changed, for example `traits` -> `identityTraits`.

**Bug Fixes:**

* Properly handle CSRF for API flows in recovery and verification strategies ([461c829](https://github.com/ory/kratos/commit/461c829dc4d7f7b70620abee2263efba78ce463a)), closes [#1141](https://github.com/ory/kratos/issues/1141)
* **session:** Use specific headers before bearer use ([82c0b54](https://github.com/ory/kratos/commit/82c0b545b29b30fcf3521d9621ec5c5f1a23dc96))
* Improve settings oas definition ([867abfc](https://github.com/ory/kratos/commit/867abfc813b08142786f71bfe28e373d4754c959))
* Use correct api spec path ([5f41f87](https://github.com/ory/kratos/commit/5f41f87bea2919cdf4e9f55c6ad938c5bc08b619))
* Use correct openapi path for validation ([#1340](https://github.com/ory/kratos/issues/1340)) ([a0f5673](https://github.com/ory/kratos/commit/a0f5673d6aa4e60bab06ef699dce231f0bf4aeff))

**Code Generation:**

* Pin v0.6.3-alpha.1 release commit ([5edf952](https://github.com/ory/kratos/commit/5edf9524d812795ac5712e4a9541b34359234724))

**Code Refactoring:**

* Improve SDK experience ([71b8511](https://github.com/ory/kratos/commit/71b8511ae1f6f77b2996a01a55accc99d171cfaf)):

    This patch resolves UX issues in the auto-generated SDKs by using consistent naming and introducing a test suite for the Ory SaaS.







**Code Generation:**

* Pin v0.6.2-alpha.1 release commit ([99c1b1d](https://github.com/ory/kratos/commit/99c1b1d674df3bd8263f7cbf1ed2bdfae6281f69))

**Documentation:**

* Update link to example email template. ([#1326](https://github.com/ory/kratos/issues/1326)) ([28a1723](https://github.com/ory/kratos/commit/28a17234b557cabf17b592ee68041aec695f6d20))






**Code Generation:**

* Pin v0.6.1-alpha.1 release commit ([1df82da](https://github.com/ory/kratos/commit/1df82daaf3f9cfd3a470d7c9bf8d96abbd52b872))

**Features:**

* Allow changing password validation API DNS name ([#1009](https://github.com/ory/kratos/issues/1009)) ([ced85e8](https://github.com/ory/kratos/commit/ced85e8091b06d864cc55c9975f8b006f6be1ce4))






**Bug Fixes:**

* Update node image ([eef307e](https://github.com/ory/kratos/commit/eef307e6bc33c9ec36ed9138f99c19f72c7be575))

**Code Generation:**

* Pin v0.6.0-alpha.2 release commit ([a3658ba](https://github.com/ory/kratos/commit/a3658badb848656b61d54b3ee35114972afc1f35))

**Features:**

* Fix unexpected emails when update profile ([#1300](https://github.com/ory/kratos/issues/1300)) ([7b24485](https://github.com/ory/kratos/commit/7b2448566f82e69d555997654ee410f9b4ff3939)), closes [#1221](https://github.com/ory/kratos/issues/1221)





## Breaking Changes


BCrypt is now the default hashing alogrithm. If you wish to continue using Argon2id please set `hashers.algorithm` to `argon2`.
This implies a significant breaking change in the verification flow payload. Please consult the new ui documentation. In essence, the login flow's `methods` key was replaced with a generic `ui` key which provides information for the UI that needs to be rendered.

To apply this patch you must apply SQL migrations. These migrations will drop the flow method table implying that all verification flows that are ongoing will become invalid. We recommend purging the flow table manually as well after this migration has been applied, if you have users doing at least one self-service flow per minute.
This implies a significant breaking change in the recovery flow payload. Please consult the new ui documentation. In essence, the login flow's `methods` key was replaced with a generic `ui` key which provides information for the UI that needs to be rendered.

To apply this patch you must apply SQL migrations. These migrations will drop the flow method table implying that all recovery flows that are ongoing will become invalid. We recommend purging the flow table manually as well after this migration has been applied, if you have users doing at least one self-service flow per minute.
This implies a significant breaking change in the settings flow payload. Please consult the new ui documentation. In essence, the login flow's `methods` key was replaced with a generic `ui` key which provides information for the UI that needs to be rendered.

To apply this patch you must apply SQL migrations. These migrations will drop the flow method table implying that all settings flows that are ongoing will become invalid. We recommend purging the flow table manually as well after this migration has been applied, if you have users doing at least one self-service flow per minute.
This implies a significant breaking change in the registration flow payload. Please consult the new ui documentation. In essence, the login flow's `methods` key was replaced with a generic `ui` key which provides information for the UI that needs to be rendered.

To apply this patch you must apply SQL migrations. These migrations will drop the flow method table implying that all registration flows that are ongoing will become invalid. We recommend purging the flow table manually as well after this migration has been applied, if you have users doing at least one self-service flow per minute.
This implies a significant breaking change in the login flow payload. Please consult the new ui documentation. In essence, the login flow's `methods` key was replaced with a generic `ui` key which provides information for the UI that needs to be rendered.

To apply this patch you must apply SQL migrations. These migrations will drop the flow method table implying that all login flows that are ongoing will become invalid. We recommend purging the flow table manually as well after this migration has been applied, if you have users doing at least one self-service flow per minute.
This change introduces a new feature: UI Nodes. Previously, all self-service flows (login, registration, ...) included form fields (e.g. `methods.password.config.fields`). However, these form fields lacked support for other types of UI elements such as links (for e.g. "Sign in with Google"), images (e.g. QR codes), javascript (e.g. WebAuthn), or text (e.g. recovery codes). With this patch, these new features have been introduced. Please be aware that this introduces significant breaking changes which you will need to adopt to in your UI. Please refer to the most recent documentation to see what has changed. Conceptionally, most things stayed the same - you do however need to update how you access and render the form fields.

Please be also aware that this patch includes SQL migrations which **purge existing self-service forms** from the database. This means that users will need to re-start the login/registration/... flow after the SQL migrations have been applied! If you wish to keep these records, make a back up of your database prior!
This change introduces a new feature: UI Nodes. Previously, all self-service flows (login, registration, ...) included form fields (e.g. `methods.password.config.fields`). However, these form fields lacked support for other types of UI elements such as links (for e.g. "Sign in with Google"), images (e.g. QR codes), javascript (e.g. WebAuthn), or text (e.g. recovery codes). With this patch, these new features have been introduced. Please be aware that this introduces significant breaking changes which you will need to adopt to in your UI. Please refer to the most recent documentation to see what has changed. Conceptionally, most things stayed the same - you do however need to update how you access and render the form fields.

Please be also aware that this patch includes SQL migrations which **purge existing self-service forms** from the database. This means that users will need to re-start the login/registration/... flow after the SQL migrations have been applied! If you wish to keep these records, make a back up of your database prior!
The configuration value for `hashers.argon2.memory` is now a string representation of the memory amount including the unit of measurement. To convert the value divide your current setting (KB) by 1024 to get a result in MB or 1048576 to get a result in GB. Example: `131072` would now become `128MB`.

Co-authored-by: aeneasr <3372410+aeneasr@users.noreply.github.com>
Co-authored-by: aeneasr <aeneas@ory.sh>
Please run SQL migrations when applying this patch.
The following configuration keys were updated:

```patch
selfservice.methods.password.config.max_breaches
```
- `password.max_breaches` -> `selfservice.methods.password.config.max_breaches`
- `password.ignore_network_errors` -> `selfservice.methods.password.config.ignore_network_errors`
After battling with [spf13/viper](https://github.com/spf13/viper) for several years we finally found a viable alternative with [knadh/koanf](https://github.com/knadh/koanf). The complete internal configuration infrastructure has changed, with several highlights:

1. Configuration sourcing works from all sources (file, env, cli flags) with validation against the configuration schema, greatly improving developer experience when changing or updating configuration.
2. Configuration reloading has improved significantly and works flawlessly on Kubernetes.
3. Performance increased dramatically, completely removing the need for a cache layer between the configuration system and ORY Hydra.
4. It is now possible to load several config files using the `--config` flag.
5. Configuration values are now sent to the tracer (e.g. Jaeger) if tracing is enabled.

Please be aware that ORY Kratos might complain about an invalid configuration, because the validation process has improved significantly.

**Bug Fixes:**

* Add include stub go files ([6d725b1](https://github.com/ory/kratos/commit/6d725b1461a26d99c8b179be8ca219ba83ba0f17))
* Add index to migration status ([8c6ec27](https://github.com/ory/kratos/commit/8c6ec2741535c090aae16f02a744f56c15923e2b))
* Add node_modules to format tasks ([e5f6b36](https://github.com/ory/kratos/commit/e5f6b36caeff080905d15566cf55f8fe4905dbc0))
* Add titles to identity schema ([73c15d2](https://github.com/ory/kratos/commit/73c15d23840aa83d2c99c013cad52ad7df285f18))
* Adopt to new go-swagger changes ([5c45bd9](https://github.com/ory/kratos/commit/5c45bd9f354bfe19b8cbcd7eb4eaebf22c441f42))
* Allow absolute file URLs as config values ([#1069](https://github.com/ory/kratos/issues/1069)) ([4bb4f67](https://github.com/ory/kratos/commit/4bb4f679d1fe0a49edb0c0189bb7a2188d4f850d))
* Allow hashtag in ui urls ([#1040](https://github.com/ory/kratos/issues/1040)) ([7591f07](https://github.com/ory/kratos/commit/7591f07f7d48376a03e9eacfdb6f4a93fd26c0d5))
* Avoid unicode-escaping ampersand in recovery URL query string ([#1212](https://github.com/ory/kratos/issues/1212)) ([d172368](https://github.com/ory/kratos/commit/d17236870af490f043d87e220179b35c9eb2dd4e))
* Bcrypt regression in credentials counting ([23fc13b](https://github.com/ory/kratos/commit/23fc13ba778e0045ca30c00d673ebd6c2f2b7fb7))
* Broken make quickstart-dev task ([#980](https://github.com/ory/kratos/issues/980)) ([999828a](https://github.com/ory/kratos/commit/999828ae036f20bde6d12fe89851e1fde9bdaca6)), closes [#965](https://github.com/ory/kratos/issues/965)
* Broken make sdk task ([#977](https://github.com/ory/kratos/issues/977)) ([5b01c7a](https://github.com/ory/kratos/commit/5b01c7a368c5bcfaa3af218d42f15288f51ab3e4)), closes [#950](https://github.com/ory/kratos/issues/950)
* Call contextualized test helpers ([e1f3f78](https://github.com/ory/kratos/commit/e1f3f7835696b039409c9d05f63665aba7a179ae))
* Code integer parsing bit size ([#1178](https://github.com/ory/kratos/issues/1178)) ([31e9632](https://github.com/ory/kratos/commit/31e9632bcd6ec3bdeabe862a4cce89021c6dd361)):

    In some cases we had a wrong bitsize of `64`, while the var was later cast to `int`. Replaced with a bitsize of `0`, which is the value to cast to `int`.

* Contextualize identity persister ([f8640c0](https://github.com/ory/kratos/commit/f8640c04f0c5873c39c8af4652d16bfbd347b79e))
* Convert all identifiers to lower case on login ([#815](https://github.com/ory/kratos/issues/815)) ([d64b575](https://github.com/ory/kratos/commit/d64b5757c710c436d6789dbdb33ed04dc11cbdf9)), closes [#814](https://github.com/ory/kratos/issues/814)
* Courier adress ([#1198](https://github.com/ory/kratos/issues/1198)) ([ebe4e64](https://github.com/ory/kratos/commit/ebe4e643150f7603a1e3a3cf6f909135097b3f49)), closes [#1194](https://github.com/ory/kratos/issues/1194)
* Courier message dequeue race condition ([#1024](https://github.com/ory/kratos/issues/1024)) ([5396a82](https://github.com/ory/kratos/commit/5396a82c34eef5d42444b5c4371bd4f820fe3eb0)), closes [#652](https://github.com/ory/kratos/issues/652) [#732](https://github.com/ory/kratos/issues/732):

    Fixes the courier message dequeuing race condition by modifying `*sql.Persister.NextMessages(ctx context.Context, limit uint8)` to retrieve only messages with status `MessageStatusQueued` and update the status of the retrieved messages to `MessageStatusProcessing` within a transaction. On message send failure, the message's status is reset to `MessageStatusQueued`, so that the message can be dequeued in a subsequent `NextMessages` call. On message send success, the status is updated to `MessageStatusSent` (no change there).

* Define credentials types as sql template and resolve crdb issue ([a2d6eeb](https://github.com/ory/kratos/commit/a2d6eeb2928c9750741237f559197fd80494310d))
* Dereference pointer types from new flow structures ([#1019](https://github.com/ory/kratos/issues/1019)) ([efedc92](https://github.com/ory/kratos/commit/efedc920e592bd6e963726e6b123ddc40df93a59))
* Do not include smtp in tracing ([#1268](https://github.com/ory/kratos/issues/1268)) ([bbfcbf9](https://github.com/ory/kratos/commit/bbfcbf9ce595d842a53a3ea21c286d5899eeb28f))
* Do not publish version at public endpoint ([3726ed4](https://github.com/ory/kratos/commit/3726ed4d145a949b25f5b5da5f58d4f448a2a90f))
* Do not reset registration method ([554bb0b](https://github.com/ory/kratos/commit/554bb0b4e62e4ac2a321fa4dbf89ffdf37b188df))
* Do not return system errors for missing identifiers ([1fcc855](https://github.com/ory/kratos/commit/1fcc8557bfee0f7ba562a635670b61dc9acb3530)), closes [#1286](https://github.com/ory/kratos/issues/1286)
* Export mailhog dockertest runner ([1384148](https://github.com/ory/kratos/commit/138414873ad319c6c32c6cc64a73547540dffc74))
* Fix random delay norm distribution math ([#1131](https://github.com/ory/kratos/issues/1131)) ([bd9d28f](https://github.com/ory/kratos/commit/bd9d28fe354710957f4ebaf71d1fffeae3968364))
* Fork audit logger from root logger ([68a09e7](https://github.com/ory/kratos/commit/68a09e7f3dc3ded9a477bb309c68ac8c4e2c2836))
* Gitlab oidc flow ([#1159](https://github.com/ory/kratos/issues/1159)) ([0bb3eb6](https://github.com/ory/kratos/commit/0bb3eb6db1144a09f4ac356cc45e1644d862bb70)), closes [#1157](https://github.com/ory/kratos/issues/1157)
* Give specific message instead of only 404 when method is disabled ([#1025](https://github.com/ory/kratos/issues/1025)) ([2f62041](https://github.com/ory/kratos/commit/2f62041a62588f5b3b062092c57053facb858e62)):

    Enabled strategies are not only used for handlers but also in other areas
    (e.g. populating the flow methods). So we should keep the logic to get
    enabled strategies and add new functions for getting all strategies.

* Ignore unset domain aliases ([ada6997](https://github.com/ory/kratos/commit/ada6997ff3dc7e48fd098e40267db5f231a5201f))
* Improve cli error output ([43e9678](https://github.com/ory/kratos/commit/43e967887280b57639565dabd92a07f02fbddeb5))
* Improve error stack trace ([4351773](https://github.com/ory/kratos/commit/43517737109088eda3b1d7f5b42f78bd5eb701d2))
* Improve error tracing ([#1005](https://github.com/ory/kratos/issues/1005)) ([456fd25](https://github.com/ory/kratos/commit/456fd254485fc80b9ae02dfca672a9fea8ae0134))
* Improve test contextualization ([2f92a70](https://github.com/ory/kratos/commit/2f92a7066d72535d32146a98207996fda45e0b96))
* Initialize randomdelay with seeded source ([9896289](https://github.com/ory/kratos/commit/9896289216f10b808a8c78b86d9c27b8d74379de))
* Insert credentials type constants as part of migrations ([#865](https://github.com/ory/kratos/issues/865)) ([92b79b8](https://github.com/ory/kratos/commit/92b79b86762edddf2ad6529b98b3383b641148d5)), closes [#861](https://github.com/ory/kratos/issues/861)
* Linking a connection may result in system error ([#990](https://github.com/ory/kratos/issues/990)) ([be02a70](https://github.com/ory/kratos/commit/be02a70c3cd60adbcc13559e1cb5dc01a8572da4)), closes [#694](https://github.com/ory/kratos/issues/694)
* Marking whoami auhorization parameter as 'in header' ([#1244](https://github.com/ory/kratos/issues/1244)) ([62d8b85](https://github.com/ory/kratos/commit/62d8b85223a0535b07620b08d35c6c3f6b127642)), closes [#1215](https://github.com/ory/kratos/issues/1215)
* Move schema loaders to correct file ([029781f](https://github.com/ory/kratos/commit/029781f69448e8abc85607a03b4bd2055158cf2c))
* Move to new transaction-safe migrations ([#1063](https://github.com/ory/kratos/issues/1063)) ([2588fb4](https://github.com/ory/kratos/commit/2588fb489d76939aeec2986d30fde9075b373831)):

    This patch introduces a new SQL transaction model for running SQL migrations. This fix is particularly targeted at CockroachDB which has limited support for mixing DDL and DML statements.

    Previously it could happen that migrations failure needed manual intervention. This has now been resolved. The new migration model is compatible with the old one and should work without a problem.

* Pass down context to registry ([0879446](https://github.com/ory/kratos/commit/08794461ed95965a9e5460ded2b4c04ab0f5e2e8))
* Re-enable SDK generation ([1d5854d](https://github.com/ory/kratos/commit/1d5854d6298e3d21f85a8fa01d3004166c4b3f50))
* Record cypress runs ([db35d8f](https://github.com/ory/kratos/commit/db35d8ff6bb44dc9e9acf131cb0a14a7f4a7d160))
* Rehydrate settings form on successful submission ([3457e1a](https://github.com/ory/kratos/commit/3457e1a46f48ed79eabff76f8af08b82f12ecc89)), closes [#1305](https://github.com/ory/kratos/issues/1305)
* Remove absolete 'make pack' from Dockerfile ([#1172](https://github.com/ory/kratos/issues/1172)) ([b8eb908](https://github.com/ory/kratos/commit/b8eb908529cc72a3147ad28e4eeee71850a8e431))
* Remove continuity cookies on errors ([85eea67](https://github.com/ory/kratos/commit/85eea6748be6ae8cdfc10cabaa6b677e4efd63eb))
* Remove include stubs ([1764e3a](https://github.com/ory/kratos/commit/1764e3a08a24db82dc391a77fdea09a91faffb5f))
* Remove obsolete clihelpers ([230fd13](https://github.com/ory/kratos/commit/230fd138d1bc7ec57647ea8eeca8e17baaacce0a))
* Remove record from bash script ([84a9315](https://github.com/ory/kratos/commit/84a9315a824cacd29d30b98b65725343af22732d))
* Remove stray non-ctx configs ([#1053](https://github.com/ory/kratos/issues/1053)) ([1fe137e](https://github.com/ory/kratos/commit/1fe137e0d6314bd0af47a29c00e2f72564e71cef))
* Remove trailing double-dot from error ([59581e3](https://github.com/ory/kratos/commit/59581e3fede0fd43028a5f064c350c3cc833b5b0))
* Remove unused sql migration ([1445d1d](https://github.com/ory/kratos/commit/1445d1d1b4b0b5e8ef3426a98ced9573063d8646))
* Remove unused var ([30a8cee](https://github.com/ory/kratos/commit/30a8cee22238d9f400e6d315a9bc99f710945f81))
* Remove verify hook ([98cfec6](https://github.com/ory/kratos/commit/98cfec6d72c2e7bf2db2e8dd6f8875e885923ba8)), closes [#1302](https://github.com/ory/kratos/issues/1302):

    The verify hook is automatically used when verification is enabled and has been removed as a configuration option.

* Replace jwt module ([#1254](https://github.com/ory/kratos/issues/1254)) ([3803c8c](https://github.com/ory/kratos/commit/3803c8ce43e35c51a9c1d7ab55bc662c398cf0d8)), closes [#1250](https://github.com/ory/kratos/issues/1250)
* Resolve build and release issues ([fb582aa](https://github.com/ory/kratos/commit/fb582aa06ad55ca3fd4e2b083e1e9bbb4ba7c715))
* Resolve clidoc issues ([599e9f7](https://github.com/ory/kratos/commit/599e9f773a743f811329cc57cea2748831105e58))
* Resolve compile issues ([63063c1](https://github.com/ory/kratos/commit/63063c15c17f4d3aca96b106275a3478a8ed717e))
* Resolve contextualized table issues ([5a4f0d9](https://github.com/ory/kratos/commit/5a4f0d92800df7fb5ca0df18203a6d73416814e1))
* Resolve crdb migration issue ([9f6edfd](https://github.com/ory/kratos/commit/9f6edfd1f544d5f85e5f5558a08672f40e928136))
* Resolve double hook invokation for registration ([032322c](https://github.com/ory/kratos/commit/032322c66fb6925d8f1473746cb4bfd800d60590))
* Resolve incorrect field types on oidc sign up completion ([f88b6ab](https://github.com/ory/kratos/commit/f88b6abe202605739092a8230fbdebaebcd4407a))
* Resolve lint issues ([0348825](https://github.com/ory/kratos/commit/03488250bcdbfda6ef6a536b4de6117fa8924dc8))
* Resolve lint issues ([75a995b](https://github.com/ory/kratos/commit/75a995b3f69778655611929b65ae22bd77c5370b))
* Resolve linting issues and disable nancy ([c8396f6](https://github.com/ory/kratos/commit/c8396f6007831240d83f77433876c5971a2191ef))
* Resolve mail queue issues ([b968bc4](https://github.com/ory/kratos/commit/b968bc4ed8962d421175adbcaa2dba6eaeea2245))
* Resolve merge regressions ([9862ac7](https://github.com/ory/kratos/commit/9862ac72e0877df4cf17c93e140c354e1ddbd0e7))
* Resolve oidc e2e regressions ([f28087a](https://github.com/ory/kratos/commit/f28087aaf133c116a81213f787dc6f2e982564c0))
* Resolve oidc regressions and e2e tests ([f5091fa](https://github.com/ory/kratos/commit/f5091fac161db0b1401b340a002278bc26891251))
* Resolve potential fsnotify leaks ([3159c0a](https://github.com/ory/kratos/commit/3159c0abe109ea4e3832770278c4e9bc4ca3b3e1))
* Resolve regressions and test failures ([8bae356](https://github.com/ory/kratos/commit/8bae3565ea5410b60c3e638a49f5454fac8e63d3))
* Resolve regressions in cookies and payloads ([9e34bf2](https://github.com/ory/kratos/commit/9e34bf2f6a2f3b007069a5415643c448798207a6))
* Resolve settings sudo regressions ([4b611f3](https://github.com/ory/kratos/commit/4b611f34755369eafcbafa2fc16da13ea3b82370))
* Resolve test regressions ([e3fb028](https://github.com/ory/kratos/commit/e3fb0281dd9be123271d11f2934cfb08fdc470b7))
* Resolve ui issues with nested form objects ([8e744b9](https://github.com/ory/kratos/commit/8e744b931954283cf5f5cbf3ebaca3fa94e035ed))
* Resolve update regression ([d0d661a](https://github.com/ory/kratos/commit/d0d661aaffcba8b039738b773c891ee6e8f6449e))
* Return delay instead of sleeping to improve tests ([27b977e](https://github.com/ory/kratos/commit/27b977ebbaa25b95caa7e3e4536a09ea0bfa61c3))
* Revert generator changes ([c18b97f](https://github.com/ory/kratos/commit/c18b97f333a638d4b4495678013c55faca4b04d0))
* Run correct error handler for registration hooks ([0d80447](https://github.com/ory/kratos/commit/0d80447102d5092e310ca728012f083147c0c5c9))
* Simplify data breaches password error reason ([#1136](https://github.com/ory/kratos/issues/1136)) ([33d29bf](https://github.com/ory/kratos/commit/33d29bf72af03aea77f1d318c19f5087a506719f)):

    This PR simplifies the error reason given when a password has appeared in data breaches to not include the actual number and rather just show "this password has appeared in data breaches and must not be used".

* Support form and json formats in decoder ([d420fe6](https://github.com/ory/kratos/commit/d420fe6e8a491b20063d4bfeaa0a841058087d32))
* Update openapi definitions for signup ([eb0b69d](https://github.com/ory/kratos/commit/eb0b69d50ce834b170186a39bbc9cda4d3366c36))
* Update quickstart node image ([c19b2f4](https://github.com/ory/kratos/commit/c19b2f4c57307e27ce289d44eff34f5aec1341da)):

    See https://github.com/ory/kratos/discussions/1301

* **cmd:** Make HTTP calls resilient ([e8ed61f](https://github.com/ory/kratos/commit/e8ed61fc3e806453f78b8fa629e96ff7b320bf95))
* **hashing:** Make bcrypt default hashing algorithm ([04abe77](https://github.com/ory/kratos/commit/04abe774ada1ef4bf318658fcf84c1d39a2a922d))
* Update to new goreleaser config ([4c2a1b7](https://github.com/ory/kratos/commit/4c2a1b7f5a0059a6e0c28779808ffb27e8910553))
* Update to new healthx ([6ec987a](https://github.com/ory/kratos/commit/6ec987ae81ef0c05f2c4d1eb836c40f9d15950b2))
* Use equalfold ([1c0e52e](https://github.com/ory/kratos/commit/1c0e52ec36ff95b53e3537c5ef457f1c818d7f6b))
* Use new TB interface ([d75a378](https://github.com/ory/kratos/commit/d75a378e700a206753f2cb17032315f2981960e7))
* Use numerical User ID instead of name to avoid k8s security warnings ([#1151](https://github.com/ory/kratos/issues/1151)) ([468a12e](https://github.com/ory/kratos/commit/468a12e56f22cfdf7bd05d68159cc735e75211b2)):

    Our docker image scanner does not allow running processes inside
    container using non-numeric User spec (to determine if we are trying
    to run docker image as root).

* Use remote dependencies ([1e56457](https://github.com/ory/kratos/commit/1e56457d49e1cde69baa41e3111ca113aa49ee3c))

**Code Generation:**

* Pin v0.6.0-alpha.1 release commit ([507d13a](https://github.com/ory/kratos/commit/507d13a8ec9cd89c9933fc8814a8a99921da69fb))

**Code Refactoring:**

* Adapt new sdk in testhelpers ([6e15f6f](https://github.com/ory/kratos/commit/6e15f6f86c0f146e846a384ffd6eac78406178bc))
* Add nid everywhere ([407fd95](https://github.com/ory/kratos/commit/407fd95889f416f0d76d6f3f43644a6fafa13b44))
* Contextualize everything ([7ebc3a9](https://github.com/ory/kratos/commit/7ebc3a9a1a2cd85d28c5a9adf2c0c8c10cbd072e)):

    This patch contextualizes all configuration and DBAL models.

* Do not use prefixed node names ([fc42ece](https://github.com/ory/kratos/commit/fc42ece24107dcb6e6a416cc54a2fb5de524fd94))
* Improve Argon2 tooling ([#961](https://github.com/ory/kratos/issues/961)) ([3151187](https://github.com/ory/kratos/commit/315118720419194be8baf5e5e64d7bf190179568)), closes [#955](https://github.com/ory/kratos/issues/955):

    This adds a load testing CLI that allows to adjust the hasher parameters under simulated load.

* Move faker to exportable module ([09f8ae5](https://github.com/ory/kratos/commit/09f8ae5755c9978574e91676bf5df6a23a2feb78))
* Move migratest helpers to ory/x ([7eca67e](https://github.com/ory/kratos/commit/7eca67eb9ec3e4ab065af7221911a74ed16c7c48))
* Move password config to selfservice ([cd0e0eb](https://github.com/ory/kratos/commit/cd0e0ebb0de372ff31c982ef023fe1979addb05a))
* Move to go 1.16 embed ([43c4a13](https://github.com/ory/kratos/commit/43c4a13c25be4a3a23a1ffdbecfaa0f9eda1a11d)):

    This patch replaces packr and pkged with the Go 1.16 embed feature.

* Remove password node attribute prefix ([e27fae4](https://github.com/ory/kratos/commit/e27fae4b0d7a91ff3964804963d4885178b80803))
* Remove profile node attribute prefix ([a3ff6f7](https://github.com/ory/kratos/commit/a3ff6f7eec45b1a9a1e7eb8569793fbc6a047d4f))
* Rename config structs and interfaces ([4a2f419](https://github.com/ory/kratos/commit/4a2f41977439354415118df3e37dd0cde8dac1aa))
* Rename form to container ([5da155a](https://github.com/ory/kratos/commit/5da155a07d3737cefabaf98c4ff650115f662480))
* Replace flow's forms with new ui node module ([647eb1e](https://github.com/ory/kratos/commit/647eb1e66850c67e539d0338cca6cb8ae476ee55))
* Replace flow's forms with new ui node module ([f74a5c2](https://github.com/ory/kratos/commit/f74a5c25af60936b59caee0866a21637a5c0ae6f))
* Replace login flow methods with ui container ([d4ca364](https://github.com/ory/kratos/commit/d4ca364fd8905cfb205ee047a9cb831064a6b9d0))
* Replace recovery flow methods with ui container ([cac0456](https://github.com/ory/kratos/commit/cac04562f2e4e77875275fcfd82c039d787607fb))
* Replace registration flow methods with ui container ([3f6388d](https://github.com/ory/kratos/commit/3f6388d03f91cfad17bd74ebca4d924b4b546668))
* Replace settings flow methods with ui container ([0efd17e](https://github.com/ory/kratos/commit/0efd17e76ba0a0cbd46916a7644b7bdf19bd4ab4))
* Replace verification flow methods with ui container ([dbf2668](https://github.com/ory/kratos/commit/dbf2668747922c93dd967961cd843354afbecfde))
* Replace viper with koanf config management ([5eb1bc0](https://github.com/ory/kratos/commit/5eb1bc0bff7c5d0f83c604484b8e845701112cad))
* Update RegisterFakes calls ([6268310](https://github.com/ory/kratos/commit/626831069ab4f971094ba0bc0b43ac9ff618d91d))
* Use underscore in webhook auth types ([26829d2](https://github.com/ory/kratos/commit/26829d21911cccd4a87c8693b6089af661c1bfe3))

**Documentation:**

* Add docker to docs main ([8ce8b78](https://github.com/ory/kratos/commit/8ce8b785e2246557253420ea97cf6b7d5ee75d58))
* Add docker to sidebar ([ed38c88](https://github.com/ory/kratos/commit/ed38c88bdbadcdcd2527a2b5270390251742bbe4))
* Add dotnet sdk ([#1183](https://github.com/ory/kratos/issues/1183)) ([32d874a](https://github.com/ory/kratos/commit/32d874a04bb384259aeb544a3fcd6b3a8b23acdd))
* Add faq sidebar ([#1105](https://github.com/ory/kratos/issues/1105)) ([10697aa](https://github.com/ory/kratos/commit/10697aa4ab5dc3e2ab90d1c037dfbe3492bf2bdf))
* Add log docs to schema config ([4967f11](https://github.com/ory/kratos/commit/4967f11d8df177ebdae855eb745e90d21ce38e9f))
* Add more HA docs ([cbb2e27](https://github.com/ory/kratos/commit/cbb2e27f8919a8991c4797a3f1c192ec364f0dd3))
* Add Rust and Dart SDKs ([6d96952](https://github.com/ory/kratos/commit/6d969528e13350ef099669510d3d37df1c007c82)):

    We now support for Rust and Dart SDKs!

* Add SameSite help ([2df6729](https://github.com/ory/kratos/commit/2df6729b4acc70532024658e8874682de64b06b3))
* Add shell-session language ([d16db87](https://github.com/ory/kratos/commit/d16db87802ae2f230a02e4deed189f473588552c))
* Add ui node docs ([e48a07d](https://github.com/ory/kratos/commit/e48a07d03c19a0677d3a56f9e57294b358f24501))
* Adding double colons ([#1187](https://github.com/ory/kratos/issues/1187)) ([fc712f4](https://github.com/ory/kratos/commit/fc712f4530066c429242491c19d1534ffb267b0c))
* Bcrypt is default and add 72 char warning ([29ae53a](https://github.com/ory/kratos/commit/29ae53a96b4472ff549b34241894d72d439c8ea1))
* Better import identities examples ([#997](https://github.com/ory/kratos/issues/997)) ([2e2880a](https://github.com/ory/kratos/commit/2e2880ac057b5c98cd69481c4f6f36b564b5871d))
* Change forum to discussions readme ([#1220](https://github.com/ory/kratos/issues/1220)) ([ae39956](https://github.com/ory/kratos/commit/ae399561ea6ed89aaadd4128bc564254984520e8))
* Describe more about Kratos login/browser flow on quickstart doc ([#1047](https://github.com/ory/kratos/issues/1047)) ([fe725ad](https://github.com/ory/kratos/commit/fe725ad12b5aed5faa8f95bec24ed3aa82512de8))
* Docker file links ([#1182](https://github.com/ory/kratos/issues/1182)) ([4d9b6a3](https://github.com/ory/kratos/commit/4d9b6a3fd5de81310016a811126e40a263ecd27c))
* Document hash timing attack mitigation ([ec86993](https://github.com/ory/kratos/commit/ec869930a9c0e6f6f56c2614835894e0a6a3eaab))
* Explain how to use `after_verification_return_to` ([7e1546b](https://github.com/ory/kratos/commit/7e1546be1fd20baca10507d642d4f209eb88dcbc))
* FAQ improvements ([#1135](https://github.com/ory/kratos/issues/1135)) ([44d0bc9](https://github.com/ory/kratos/commit/44d0bc968a7c0ba5c0793b2349820fa8133bada3))
* FAQ item & minor changes ([#1174](https://github.com/ory/kratos/issues/1174)) ([11cf630](https://github.com/ory/kratos/commit/11cf630082b56c80d12f5915f8e34aa03a7e8c54))
* Fix broken link ([#1037](https://github.com/ory/kratos/issues/1037)) ([6b9aae8](https://github.com/ory/kratos/commit/6b9aae8af5aa3bd614c99b32e341fbd533caf116))
* Fix failing build ([0de328f](https://github.com/ory/kratos/commit/0de328ff0053605e6bded589a79d3ab938d55b31))
* Fix formatting ([#966](https://github.com/ory/kratos/issues/966)) ([687251a](https://github.com/ory/kratos/commit/687251a24e796322b43f8aed6b1fb3d7900e3271))
* Fix identity state bullets ([#1095](https://github.com/ory/kratos/issues/1095)) ([f476334](https://github.com/ory/kratos/commit/f476334c4693277656ad88e768f66b59cbcba126))
* Fix known/unknown email account recovery ([#1211](https://github.com/ory/kratos/issues/1211)) ([e208ca5](https://github.com/ory/kratos/commit/e208ca50ba4f03d5410c9644aaa3b04bdf1b8dbd))
* Fix link ([7f6d7f5](https://github.com/ory/kratos/commit/7f6d7f501d7118dfe6868c9d923fb5ecc5eded48))
* Fix link ([#1128](https://github.com/ory/kratos/issues/1128)) ([e7043e9](https://github.com/ory/kratos/commit/e7043e9b99260eaff2b48ca6f457af46a1521654))
* Fix link to blogpost ([#949](https://github.com/ory/kratos/issues/949)) ([4622e32](https://github.com/ory/kratos/commit/4622e3228fb12231222c7e6b602458111f35f727)), closes [#945](https://github.com/ory/kratos/issues/945)
* Fix link to self-service flows overview ([#995](https://github.com/ory/kratos/issues/995)) ([2be8778](https://github.com/ory/kratos/commit/2be877847644a3df2645ac3be4bbd7704db30b17))
* Fix note block in third party login guide ([#920](https://github.com/ory/kratos/issues/920)) ([745cea0](https://github.com/ory/kratos/commit/745cea02d0e9940f689e668bbd814b29fd53bf37)):

    Allows the document to render properly

* Fix npm links ([#991](https://github.com/ory/kratos/issues/991)) ([4ce4468](https://github.com/ory/kratos/commit/4ce4468132dde21c1692e3a834ad7780bee12b90))
* Fix self-service code flows labels ([#1253](https://github.com/ory/kratos/issues/1253)) ([f2ed424](https://github.com/ory/kratos/commit/f2ed424289cdd2a0edc1736888dd15be6df65f11))
* Fix typo in README ([#1122](https://github.com/ory/kratos/issues/1122)) ([e500707](https://github.com/ory/kratos/commit/e5007078c3cd597cea669827b96c7e6f205f2f32))
* Link to argon2 blogpost and add cross-references ([#1038](https://github.com/ory/kratos/issues/1038)) ([9ab7c3d](https://github.com/ory/kratos/commit/9ab7c3df59ecd94a74a7bf18af9c0ded5305e042))
* Make explicit the ID of the default schema ([#1173](https://github.com/ory/kratos/issues/1173)) ([cc6e9ff](https://github.com/ory/kratos/commit/cc6e9ffbac7118436d85078720cde2de98a68044))
* Minor cosmetics ([#1050](https://github.com/ory/kratos/issues/1050)) ([34db06f](https://github.com/ory/kratos/commit/34db06fd4f83d415c09109b06dfd3b82ce03705e))
* Minor improvements ([#1052](https://github.com/ory/kratos/issues/1052)) ([f0672b5](https://github.com/ory/kratos/commit/f0672b5cb8cca41fa914db21798d20f00a5699f9))
* ORY -> Ory ([ea30979](https://github.com/ory/kratos/commit/ea309797bf59f3da5c5cd184e45f2e585144be56))
* Reformat settings code samples ([cdbbf4d](https://github.com/ory/kratos/commit/cdbbf4df5fa3fa667a78d5cf682bc7fa36693e9d))
* Remove unnecessary and wrong docker pull commands ([#1203](https://github.com/ory/kratos/issues/1203)) ([2b0342a](https://github.com/ory/kratos/commit/2b0342ad7607d705bcebfafd5a78e4e09e57a940))
* Resolve duplication error ([a3d8284](https://github.com/ory/kratos/commit/a3d8284ab20ae76bccba361601b7290af20bdde6))
* Update build from source ([9b5754f](https://github.com/ory/kratos/commit/9b5754f36661f6de9c95f30c06f28164fe5be48b)), closes [#979](https://github.com/ory/kratos/issues/979)
* Update email template docs ([1778cb9](https://github.com/ory/kratos/commit/1778cb9a293feb2c91c0b1921ab78a0395cdca98)), closes [#897](https://github.com/ory/kratos/issues/897)
* Update identity-data-model links ([b5fd9a3](https://github.com/ory/kratos/commit/b5fd9a3a0821215f94da168c9c6f87dceba8c8f4))
* Update identity.ID field documentation ([4624f03](https://github.com/ory/kratos/commit/4624f03a5e9249a5449992a1f0b7ec80dc3499fd)):

    See https://github.com/ory/kratos/discussions/956

* Update kratos video link ([#1073](https://github.com/ory/kratos/issues/1073)) ([e86178f](https://github.com/ory/kratos/commit/e86178f4ee66e5053e0da2fab2c21ecb2e730ada))
* Update login code samples ([695a30f](https://github.com/ory/kratos/commit/695a30f6c80f277676bf04b4665efeb7ea4db618))
* Update login code samples ([ce6c755](https://github.com/ory/kratos/commit/ce6c75587bea80ef83855d764fed79a9d6c948d3))
* Update quickstart samples ([c3fcaba](https://github.com/ory/kratos/commit/c3fcaba65899d9d46a08ca8b60ec0c010f70b16c))
* Update recovery code samples ([d9fbb62](https://github.com/ory/kratos/commit/d9fbb62faff5144f587136935f15d24b6399f29c))
* Update registration code samples ([317810f](https://github.com/ory/kratos/commit/317810ffd8ba6faf87f2248263b6c82cf4e9ffd8))
* Update self-service code samples ([6415011](https://github.com/ory/kratos/commit/6415011ab83a19972c6f52467055fbdcef23a0cc))
* Update settings code samples ([bbd6266](https://github.com/ory/kratos/commit/bbd6266c22097fae195654957cbab589d04892c7))
* Update verification code samples ([4285dec](https://github.com/ory/kratos/commit/4285dec59a8fc31fa3416b594c765f5da9a9de1c))
* Use correct extension for identity-data-model ([acab3e8](https://github.com/ory/kratos/commit/acab3e8b489d9865e4bf0805895f0b7ae9e6f1b8)), closes [/github.com/ory/kratos/pull/1197#issuecomment-819455322](https://github.com//github.com/ory/kratos/pull/1197/issues/issuecomment-819455322)
* **prometheus:** Update codedoc ([47146ea](https://github.com/ory/kratos/commit/47146ea8ce169ee908aa4d33b59a01e9df4bae10))

**Features:**

* Add email template specification in doc ([#898](https://github.com/ory/kratos/issues/898)) ([4230d9e](https://github.com/ory/kratos/commit/4230d9e0fc35c651b0d2cbdbbf9e1f1c514743f8))
* Add error for when no login strategy was found ([6bae66c](https://github.com/ory/kratos/commit/6bae66cde362c4e2995c9d06a0d3ffee403feb74))
* Add facebook provider to oidc providers and documentation ([#1035](https://github.com/ory/kratos/issues/1035)) ([905bb03](https://github.com/ory/kratos/commit/905bb032520189212bd88f29641903945ae03608)), closes [#1034](https://github.com/ory/kratos/issues/1034)
* Add FAQ to docs ([#1096](https://github.com/ory/kratos/issues/1096)) ([9c6b68c](https://github.com/ory/kratos/commit/9c6b68c454f472b26c34e1975b6a67b24b218f47))
* Add gh login to claims ([49deb2e](https://github.com/ory/kratos/commit/49deb2e166362a5d051bc08523ef44425f144bdd))
* Add login strategy text message ([7468c83](https://github.com/ory/kratos/commit/7468c835d4800c207035897fc9962860d8ab7803))
* Add more tests for multi domain args ([e99803b](https://github.com/ory/kratos/commit/e99803b62a847bcee52bcd87fa8088124b4deae2))
* Add Prometheus monitoring to Public APIs ([#1022](https://github.com/ory/kratos/issues/1022)) ([75a4f1a](https://github.com/ory/kratos/commit/75a4f1a5472ffd780fed43a7395a191ed495c6e9))
* Add random delay to login flow ([#1088](https://github.com/ory/kratos/issues/1088)) ([cb9894f](https://github.com/ory/kratos/commit/cb9894fefc694a4092215d3981e80f287021542f)), closes [#832](https://github.com/ory/kratos/issues/832)
* Add return_url to verification flow ([#1149](https://github.com/ory/kratos/issues/1149)) ([bb99912](https://github.com/ory/kratos/commit/bb99912d823e9bcffa41edf50a01dcae40117fe6)), closes [#1123](https://github.com/ory/kratos/issues/1123) [#1133](https://github.com/ory/kratos/issues/1133)
* Add sql migrations for new login flow ([e947edf](https://github.com/ory/kratos/commit/e947edf497b36bc576061c9ae38049e84ee48575))
* Add sql tracing ([3c4cc1c](https://github.com/ory/kratos/commit/3c4cc1cec170df14331288170a94ada770d3289f))
* Add tracing to config schema ([007dde4](https://github.com/ory/kratos/commit/007dde4482d11f22b8527c94b002da675152a872))
* Add transporter with host modification ([2c41b81](https://github.com/ory/kratos/commit/2c41b81be947f9972638d082105f0f5c83078b91))
* Add workaround template for go openapi ([5d72d10](https://github.com/ory/kratos/commit/5d72d10f6c6948c48c5701fe348084a668c8311a))
* Adds slack sogial login ([#974](https://github.com/ory/kratos/issues/974)) ([7c66053](https://github.com/ory/kratos/commit/7c66053390b3086fe7233625038a78431a61e507)), closes [#953](https://github.com/ory/kratos/issues/953)
* Allow session cookie name configuration ([77ce316](https://github.com/ory/kratos/commit/77ce3162ba97cf5c516c26ef499d9fa892162f0a)), closes [#268](https://github.com/ory/kratos/issues/268)
* Allow specifying sender name in smtp.from_address ([#1100](https://github.com/ory/kratos/issues/1100)) ([5904fe3](https://github.com/ory/kratos/commit/5904fe319f75f8138783434d568db6fc7c55b301))
* Bcrypt algorithm support ([#1169](https://github.com/ory/kratos/issues/1169)) ([b2612ee](https://github.com/ory/kratos/commit/b2612eefbad98d29482d364f670549f470d0a6f5)):

    This patch adds the ability to use BCrypt instead of Argon2id for password hashing. We recommend using BCrypt for web workloads where password hashing should take around 200ms. For workloads where login takes >= 2 seconds, we recommend to continue using Argon2id.

    To use bcrypt for password hashing, set your config as follows:

     ```
    hashers:
     bcrypt:
        cost: 12
      algorithm: bcrypt
     ```

    Switching the hashing algorithm will not break existing passwords!


    Co-authored-by: Patrik <zepatrik@users.noreply.github.com>

* Check migrations in health check ([c6ef7ad](https://github.com/ory/kratos/commit/c6ef7ad16b70310c645550f7e41b3c8aff847de3))
* Configure domain alias as query param ([9d8563e](https://github.com/ory/kratos/commit/9d8563eeb3293c42cce440ad74f025b304cccbbe))
* Contextualize configuration ([d3d5327](https://github.com/ory/kratos/commit/d3d5327a3622318265a063be4782caa25e645a05))
* Contextualize health checks ([8145a1c](https://github.com/ory/kratos/commit/8145a1c9acaeab441e787118d40ccd448ea82fe4))
* Contextualize http client in cli calls ([3b3ef8f](https://github.com/ory/kratos/commit/3b3ef8f025d75b244d9285036e66f79af7d5ee35))
* Contextualize persitence testers ([6440373](https://github.com/ory/kratos/commit/64403736ad9f8b264567e1f8eed1af710cab6046))
* Courier foreground worker with "kratos courier watch" ([#1062](https://github.com/ory/kratos/issues/1062)) ([500b8ba](https://github.com/ory/kratos/commit/500b8bacd9fd541afd053f42fec66443cfebabda)), closes [#1033](https://github.com/ory/kratos/issues/1033) [#1024](https://github.com/ory/kratos/issues/1024):

    BREACKING CHANGES: This patch moves the courier watcher (responsible for sending mail) to its own foreground worker, which can be executed as a, for example, Kubernetes job.

    It is still possible to have the previous behaviour which would run the worker as a background task when running `kratos serve` by using the `--watch-courier` flag.

    To run the foreground worker, use `kratos courier watch -c your/config.yaml`.

* Do not enforce bcrypt 12 for dev envs ([bbf44d8](https://github.com/ory/kratos/commit/bbf44d887ae5cdb5975516149c74b3ba10896209))
* Email input validation ([#1287](https://github.com/ory/kratos/issues/1287)) ([cd56b73](https://github.com/ory/kratos/commit/cd56b73df363dd37485f07d31fef11fd4d9f40a6)), closes [#1285](https://github.com/ory/kratos/issues/1285)
* Export and add config options ([4391fe5](https://github.com/ory/kratos/commit/4391fe572eb6a766afe9808396847ca5fdca07f5))
* Expose courier worker ([f50969e](https://github.com/ory/kratos/commit/f50969ecba757dea558e9e8b9dd142f5f564d53a))
* Expose crdb ui ([504d518](https://github.com/ory/kratos/commit/504d5181f5e391bb8d67768b314a0348ed252c8b))
* Global docs sidebar ([#1258](https://github.com/ory/kratos/issues/1258)) ([7108262](https://github.com/ory/kratos/commit/71082624e093b8c100e71ae59050f89b35ac20a2))
* Implement and test domain aliasing ([1516a54](https://github.com/ory/kratos/commit/1516a54657df485627251de4e7019bc16353c956)):

    This patch adds a feature called domain aliasing. For more information, head over to http://ory.sh/docs/kratos/next/guides/multi-domain-cookies

* Improve oas spec and fix mobile tests ([4ead2c8](https://github.com/ory/kratos/commit/4ead2c826a2f1a307e327b9736dd8ac99ef52743))
* Improve sorting of ui fields ([797b49d](https://github.com/ory/kratos/commit/797b49d0175280f85f568014cf3083e9bc42d354)):

    See https://github.com/ory/kratos/discussions/1196

* Include schema ([348a493](https://github.com/ory/kratos/commit/348a493c9e5381830b76e57cad803a308e6ce53a))
* Make cli commands consumable in Ory Cloud ([#926](https://github.com/ory/kratos/issues/926)) ([fed790b](https://github.com/ory/kratos/commit/fed790b0f71f028f6d92e8ebceee188dbdb20770))
* Migrate to openapi v3 ([595224b](https://github.com/ory/kratos/commit/595224b1efd5a225702ef236a87f08180a7118b8))
* Populate email templates at delivery time, add plaintext defaults ([#1155](https://github.com/ory/kratos/issues/1155)) ([7749c7a](https://github.com/ory/kratos/commit/7749c7a75a4386c1fd53db57626355467b698c2f)), closes [#1065](https://github.com/ory/kratos/issues/1065)
* Sort and label nodes with easy to use defaults ([cbec27c](https://github.com/ory/kratos/commit/cbec27c957a733411e4c1d511ed5854855b7236e)):

    Ory Kratos takes a guess based on best practices for

    - ordering UI nodes (e.g. email, password, submit button)
    - grouping UI nodes (e.g. keep password and oidc nodes together)
    - labeling UI nodes (e.g. "Sign in with GitHub")
    - using the "title" attribute from the identity schema to label trait fields

    This greatly simplifies front-end code on your end and makes it even easier to integrate with Ory Kratos! If you want a custom experience with e.g. translations or other things you can always adjust this in your UI integration!

* Support base64 inline schemas ([815a248](https://github.com/ory/kratos/commit/815a24890a118f4128ac083241a93d8df27042f7))
* Support contextual csrf cookies ([957ef38](https://github.com/ory/kratos/commit/957ef38b69fc6ab071b91262736e6c191be3a4b8))
* Support domain aliasing in session cookie ([0681c12](https://github.com/ory/kratos/commit/0681c123f2d856ca27caee645dadc9e6e3731d2c))
* Support label in oidc config ([a99cdcd](https://github.com/ory/kratos/commit/a99cdcddaa0c4bd7b679884b232c2ef8f2dcd978))
* Support retryable CRDB transactions ([f0c21d7](https://github.com/ory/kratos/commit/f0c21d7e0a6ed85818d0e9025a451cb8cbdee086))
* Unix sockets support ([#1255](https://github.com/ory/kratos/issues/1255)) ([ad010de](https://github.com/ory/kratos/commit/ad010de240ddd9219f0cfb2ca3fbb180d2d3a697))
* Web hooks support (recovery) ([#1289](https://github.com/ory/kratos/issues/1289)) ([3e181fe](https://github.com/ory/kratos/commit/3e181fe3d7750a715ab31eb8347fbb4bdb89d6e6)), closes [#271](https://github.com/ory/kratos/issues/271):

    feat: web hooks for self-service flows

    This feature adds the ability to define web-hooks using a mixture of configuration and JsonNet. This allows integration with services like Mailchimp, Stripe, CRMs, and all other APIs that support REST requests. Additional to these new changes it is now possible to define hooks for verification and recovery as well!

    For more information, head over to the [hooks documentation](https://www.ory.sh/kratos/docs/self-service/hooks).

* **courier:** Allow sending individual messages ([cbb2c0b](https://github.com/ory/kratos/commit/cbb2c0bef63323a177589e9d2a809c84b4f1acdd))
* **oidc:** Support google hd claim ([#1097](https://github.com/ory/kratos/issues/1097)) ([1f20a5c](https://github.com/ory/kratos/commit/1f20a5ceba7682719112d24a3b18bf046fb2ac22))
* **schema:** Add totp errors ([a61f881](https://github.com/ory/kratos/commit/a61f8814101401dbb422967e37b6c6c1ae85d113))

**Tests:**

* Add case to ensure correct behavior when verifying a different email address ([#999](https://github.com/ory/kratos/issues/999)) ([f95a117](https://github.com/ory/kratos/commit/f95a117677c9c59436ad10aa8951fe875c39a64f)), closes [#998](https://github.com/ory/kratos/issues/998)
* Add oasis test case ([f80691b](https://github.com/ory/kratos/commit/f80691b9dd77566857c4284e2639cc94d5b8c333))
* Bump poll interval ([b3dc925](https://github.com/ory/kratos/commit/b3dc925a5d43557293745ee81c0ffb3db37b6342))
* Bump video quality ([b7f8d04](https://github.com/ory/kratos/commit/b7f8d042646037e1589ae2d03602bd63a5cec2fe))
* Bump wait times ([b2e43f8](https://github.com/ory/kratos/commit/b2e43f8b0b64784f60e5f57d9a0f5d2928c2b891))
* Clean up hydra env before restart ([cf49414](https://github.com/ory/kratos/commit/cf494149e6a46b15e3b174185e1e87cfcd6f9f7a))
* Longer wait times ([4bec9ef](https://github.com/ory/kratos/commit/4bec9ef50f14f22342a311f09ba1b59cde47befc))
* Reliable migration tests on crdb ([2e3764b](https://github.com/ory/kratos/commit/2e3764ba66c156d810de66fba2b0e142dced6f4d))
* Remove old noop test ([16dca3f](https://github.com/ory/kratos/commit/16dca3f78b2021c09ec83e81ab6d2e68c42ca081))
* Resolve compile issues ([c1b5ba4](https://github.com/ory/kratos/commit/c1b5ba42171ec522579df9dfaff27b5b74a1566a))
* Resolve flaky tests ([cb670a8](https://github.com/ory/kratos/commit/cb670a854cbb09b8437bfed7e4a6908ff6dcfd27))
* Resolve json parser test regression ([a1b9b9a](https://github.com/ory/kratos/commit/a1b9b9a95d58583dc7ecf6d2a501da52f84dd6bb))
* Resolve login integration regressions ([388b5b2](https://github.com/ory/kratos/commit/388b5b27d6dee7770e5f37d6d83c532044a4e984))
* Resolve migration regression ([2051a71](https://github.com/ory/kratos/commit/2051a716cb4b8cf334dd65f2ccddb31e5fbed545))
* Resolve more json parser test regressions ([ff791c4](https://github.com/ory/kratos/commit/ff791c41a1d9ce25af4e883469d3f8c0ef9eb302))
* Resolve regression ([e2b0ad3](https://github.com/ory/kratos/commit/e2b0ad3c1845da80f078b11b327b9a0376cbb7c5))
* Update schema tests for webhooks ([d1ddfa8](https://github.com/ory/kratos/commit/d1ddfa80742728b28dc5710ca5b6e7282a2dec55))
* **e2e:** Significantly reduce wait and idle times ([f525fc5](https://github.com/ory/kratos/commit/f525fc53afec6f5232ce507fe25ddec1b9069196))
* Resolve more regressions ([c5a23af](https://github.com/ory/kratos/commit/c5a23af81427480088651833d904e3403a969fab))
* Resolve order regression ([40a849c](https://github.com/ory/kratos/commit/40a849ca35f4700185322e9ac4f6a4b70132851c))
* Resolve regression ([f0c9e5f](https://github.com/ory/kratos/commit/f0c9e5ff105d76d6bc9478c98522b2440c7181df))
* Resolve regressions ([4b9da3c](https://github.com/ory/kratos/commit/4b9da3c9d98d40f7b71a56c51543fc115974630d))
* Resolve stub regressions ([82650cf](https://github.com/ory/kratos/commit/82650cf1843f6bfde015f556f4452a7b6fd52b11))
* Resolve test migrations ([de0b65d](https://github.com/ory/kratos/commit/de0b65d96daef0e31c12b3b6915f283a8e71244b))
* Resolve test regression issues ([ccf9fed](https://github.com/ory/kratos/commit/ccf9feddade11f9fcaaf1c37dd3efeb2c4df6649))
* Speed up tests ([a16737c](https://github.com/ory/kratos/commit/a16737cccc36a14444711660f1737913ffd7ba01))
* Update test description ([55fb37f](https://github.com/ory/kratos/commit/55fb37f62fc3ab7c0d5324ed31ef3e7f66a73aa2))
* Use bcrypt cost 4 to reduce CI times ([cabe97d](https://github.com/ory/kratos/commit/cabe97d0656858fd1ee0442b40881417e91294f3))
* Use fast bcrypt for e2e ([d90cf13](https://github.com/ory/kratos/commit/d90cf13230632e76eb74965c0945573b4f2e98ff))

**Unclassified:**

* Format ([e4b7e79](https://github.com/ory/kratos/commit/e4b7e79f4ee91dadfcd008a5b3e318b6bfedad10))
* Format ([193d266](https://github.com/ory/kratos/commit/193d2668ae0955a1346390057539a8b796d17afd))
* Format ([1ebfbde](https://github.com/ory/kratos/commit/1ebfbdea75f27c8eeafa7d3aff45de133ea340bb))
* Format ([ba1eeef](https://github.com/ory/kratos/commit/ba1eeef4f232c4ab59343a2ca3c7cf0eb6dfd110))
* Format ([ada5dbb](https://github.com/ory/kratos/commit/ada5dbb58c45502b8275850a3bc0876debc66888))
* Initial documentation tests via Text-Runner ([#567](https://github.com/ory/kratos/issues/567)) ([c30eb26](https://github.com/ory/kratos/commit/c30eb26f76ab70a6098c0b40c9a04726d36d72f2))
*  fix: resolve clidoc issues (#976) ([346bc73](https://github.com/ory/kratos/commit/346bc73921655d52861b8803eb3351c4205657ee)), closes [#976](https://github.com/ory/kratos/issues/976) [#951](https://github.com/ory/kratos/issues/951)
* Format ([17a0bf5](https://github.com/ory/kratos/commit/17a0bf5872b33eac615afc675c7d92d7c7441b2e))
* :bug: fix ory home directory path (#897) ([2fca2be](https://github.com/ory/kratos/commit/2fca2bedaa907691bef324c11545e007b51d4881)), closes [#897](https://github.com/ory/kratos/issues/897)
* Fix typo in config schema ([16337f1](https://github.com/ory/kratos/commit/16337f13e4388a715c8109c29cf198c82a848a16))






**Bug Fixes:**

* CSRF token is required when using the Revoke Session API endpoint ([#839](https://github.com/ory/kratos/issues/839)) ([d3218a0](https://github.com/ory/kratos/commit/d3218a0f23de7293b0a4a966ad21369a92b68b1a)), closes [#838](https://github.com/ory/kratos/issues/838)
* Incorrect home path ([#848](https://github.com/ory/kratos/issues/848)) ([5265af0](https://github.com/ory/kratos/commit/5265af00c92fe505819300caddfcc64004d45c65))
* Make password policy configurable ([#888](https://github.com/ory/kratos/issues/888)) ([7a00483](https://github.com/ory/kratos/commit/7a00483908bb623efdf281e76005c4485ea6b1ab)), closes [#450](https://github.com/ory/kratos/issues/450) [#316](https://github.com/ory/kratos/issues/316):

    Allows configuring password breach thresholds and optionally enforces checks against the HIBP API.

* Remove obsolete types ([#887](https://github.com/ory/kratos/issues/887)) ([b8bac7a](https://github.com/ory/kratos/commit/b8bac7aa56c16cd98f76a95a5e0d01fb1bbde6b7)), closes [#716](https://github.com/ory/kratos/issues/716)
* Set samesite attribute to lax if in dev mode ([#824](https://github.com/ory/kratos/issues/824)) ([91d6698](https://github.com/ory/kratos/commit/91d6698e4ce05ee59bb72fc84b54af9d1d204b41)), closes [#821](https://github.com/ory/kratos/issues/821)
* Use working cache-control header for cdn/proxies/cache ([#869](https://github.com/ory/kratos/issues/869)) ([d8e3d40](https://github.com/ory/kratos/commit/d8e3d40001ffdc64da2288f3cffd53cf3bfdf781)), closes [#601](https://github.com/ory/kratos/issues/601)

**Code Generation:**

* Pin v0.5.5-alpha.1 release commit ([83aedcb](https://github.com/ory/kratos/commit/83aedcb885acb96c5deb39fff675d5f0528af32d))

**Documentation:**

* Add contributing to sidebar ([#866](https://github.com/ory/kratos/issues/866)) ([44f33f9](https://github.com/ory/kratos/commit/44f33f97d43f2a3c553a65ebb2986e0731c0e5f2)):

    The same change as in https://github.com/ory/hydra/pull/2209

* Add newsletter to config ([1735ca2](https://github.com/ory/kratos/commit/1735ca2ced104971de4e97524d0a23d57ba045f2))
* Add recovery flow  ([#868](https://github.com/ory/kratos/issues/868)) ([d95cfe9](https://github.com/ory/kratos/commit/d95cfe9759d3ffc08c24048a064c0c800abdf4b4)), closes [#864](https://github.com/ory/kratos/issues/864):

    Added a short section for the recovery flow on managing-user-identities.

* Fix account recovery click instruction ([#870](https://github.com/ory/kratos/issues/870)) ([383de9e](https://github.com/ory/kratos/commit/383de9ecf6f6504dbb9c20fb4cb984e934f0751e))
* Fix broken link ([#893](https://github.com/ory/kratos/issues/893)) ([dec38a2](https://github.com/ory/kratos/commit/dec38a28964aaa13827d356e5bfa12c2a6d1400e)), closes [#835](https://github.com/ory/kratos/issues/835)
* Fix oidc config example structure ([#845](https://github.com/ory/kratos/issues/845)) ([c102a68](https://github.com/ory/kratos/commit/c102a6844db29f994b67d23bb04e64ee71376264))
* Fix redirect ([#802](https://github.com/ory/kratos/issues/802)) ([b868782](https://github.com/ory/kratos/commit/b86878229f343e6b11521596b04040f892d1e2c3))
* Fix typo ([#847](https://github.com/ory/kratos/issues/847)) ([9b3da9f](https://github.com/ory/kratos/commit/9b3da9f0fe2ce71743115844d8c91a1dc9c4cbae))
* Fix typo ([#881](https://github.com/ory/kratos/issues/881)) ([3078293](https://github.com/ory/kratos/commit/3078293717a2ce21c4b939de4c2c4886c75303b5))
* Fix typo MKFA to MFA ([#826](https://github.com/ory/kratos/issues/826)) ([a5613d0](https://github.com/ory/kratos/commit/a5613d08aa21f90f4d192e5663ba4977b3de16c3))
* Remove workaround note ([#886](https://github.com/ory/kratos/issues/886)) ([05409bc](https://github.com/ory/kratos/commit/05409bc13f527398e3de01f29437e5d4353ef8d4)), closes [#718](https://github.com/ory/kratos/issues/718)
* Swagger specs for selfservice settings browser flow ([#825](https://github.com/ory/kratos/issues/825)) ([28d50f4](https://github.com/ory/kratos/commit/28d50f45ab14d561609be7047cac13902394b547))
* Update oidc provider with json conf support ([#833](https://github.com/ory/kratos/issues/833)) ([670eb37](https://github.com/ory/kratos/commit/670eb37d19674f33a36402cd9a88d61ca7327751))

**Features:**

* Add return_to parameter to logout flow ([#823](https://github.com/ory/kratos/issues/823)) ([1c146dd](https://github.com/ory/kratos/commit/1c146dd21d616a56f510019abadd37402782bb39)), closes [#702](https://github.com/ory/kratos/issues/702)
* Add selinux compatible quickstart config ([#889](https://github.com/ory/kratos/issues/889)) ([0f87948](https://github.com/ory/kratos/commit/0f879481df209ed96b778799adcc2a9424449b37)), closes [#831](https://github.com/ory/kratos/issues/831)

**Tests:**

* Ensure registration runs only once ([#872](https://github.com/ory/kratos/issues/872)) ([5ffc036](https://github.com/ory/kratos/commit/5ffc036ac82f36ad6ef499e217971275a35fc23a))

**Unclassified:**

* docs: fix link and typo in Configuring Cookies (#883) ([c51ed6b](https://github.com/ory/kratos/commit/c51ed6b789d2e3a8fe4e93565c3bded37d298f98)), closes [#883](https://github.com/ory/kratos/issues/883)






**Bug Fixes:**

* Case in settings handler method ([#798](https://github.com/ory/kratos/issues/798)) ([83eb4e0](https://github.com/ory/kratos/commit/83eb4e0021621014d2b543e57a01401381f07fe4))
* Force brew install statement ([#796](https://github.com/ory/kratos/issues/796)) ([ad542ad](https://github.com/ory/kratos/commit/ad542ad5919205ac26a757145474e5a46f3937ec)):

    Closes https://github.com/ory/homebrew-kratos/issues/1


**Code Generation:**

* Pin v0.5.4-alpha.1 release commit ([b02926c](https://github.com/ory/kratos/commit/b02926c42aee2748bc37ce2600596bd0c2537a0d))

**Code Refactoring:**

* Move pkger and ioutil helpers to ory/x ([60a0fc4](https://github.com/ory/kratos/commit/60a0fc449d90ead6065ca00926536a989d8b2a2b))

**Documentation:**

* Fix another broken link ([15bae9f](https://github.com/ory/kratos/commit/15bae9f893c2e2910167326d987455246c110001))
* Fix broken links ([#795](https://github.com/ory/kratos/issues/795)) ([0ab0e7e](https://github.com/ory/kratos/commit/0ab0e7eca8e95d6c26d028c177cbbd1f06b68871)), closes [#793](https://github.com/ory/kratos/issues/793)
* Fix broken relative link ([#812](https://github.com/ory/kratos/issues/812)) ([b32b173](https://github.com/ory/kratos/commit/b32b173fe30b7c5c43700abfa4ddb3409a33556b))
* Fix links ([#800](https://github.com/ory/kratos/issues/800)) ([5fcc272](https://github.com/ory/kratos/commit/5fcc272e625de9e583b2ec24d5679895a6d24c1b))
* Fix oidc config examples ([#799](https://github.com/ory/kratos/issues/799)) ([8a4f480](https://github.com/ory/kratos/commit/8a4f480121995d9899668f037382086fcdd2da4c))
* Fix self-service recovery flow typo ([#807](https://github.com/ory/kratos/issues/807)) ([800110d](https://github.com/ory/kratos/commit/800110d87c9df70a5ec79b58d9fcb9ae39ff76b9))
* Remove duplicate words & fix spelling ([#810](https://github.com/ory/kratos/issues/810)) ([4e1b966](https://github.com/ory/kratos/commit/4e1b96667d9f08dbafeb2f5ce144ca43309de8e0))
* Remove leftover category from reference sidebar ([#813](https://github.com/ory/kratos/issues/813)) ([94fde51](https://github.com/ory/kratos/commit/94fde5101d00b9e1f7228e9d122ef0a8e4719355))
* Use correct links ([#797](https://github.com/ory/kratos/issues/797)) ([a4de293](https://github.com/ory/kratos/commit/a4de29399e4f1b5d0a33acc85478f2d38579a174))

**Features:**

* Add helper for choosing argon2 parameters ([#803](https://github.com/ory/kratos/issues/803)) ([ca5a69b](https://github.com/ory/kratos/commit/ca5a69b798635d0e5361fd5b0cc369b035dca738)), closes [#723](https://github.com/ory/kratos/issues/723) [#572](https://github.com/ory/kratos/issues/572) [#647](https://github.com/ory/kratos/issues/647):

    This patch adds the new command "hashers argon2 calibrate" which allows one to pick the desired hashing time for password hashing and then chooses the optimal parameters for the hardware the command is running on:

    ```
    $ kratos hashers argon2 calibrate 500ms
    Increasing memory to get over 500ms:
        took 2.846592732s in try 0
        took 6.006488824s in try 1
      took 4.42657975s with 4.00GB of memory
    [...]
    Decreasing iterations to get under 500ms:
        took 484.257775ms in try 0
        took 488.784192ms in try 1
      took 486.534204ms with 3 iterations
    Settled on 3 iterations.

    {
      "memory": 1048576,
      "iterations": 3,
      "parallelism": 32,
      "salt_length": 16,
      "key_length": 32
    }
    ```







**Bug Fixes:**

* Add "x-session-token" to default allowed headers ([3c912e4](https://github.com/ory/kratos/commit/3c912e4c7d46fd45c00cabb68ed7770bd44f7d07))
* Do not set cookies on api endpoints ([2f67c28](https://github.com/ory/kratos/commit/2f67c28718856ea03ea2effa89b28a8c4b3b8ae0))
* Do not set csrf cookies on potential api endpoints ([4d97a95](https://github.com/ory/kratos/commit/4d97a95d084ea99f5aca158609e197acd256cdd7))
* Ignore unsupported migration dialects ([12bb8d1](https://github.com/ory/kratos/commit/12bb8d14ae1edef18591996411be67d5693e5101)), closes [#778](https://github.com/ory/kratos/issues/778):

    Skips sqlite3 migrations when support is lacking.

* Improve semver regex ([584c0b5](https://github.com/ory/kratos/commit/584c0b5043e85e88ac2648cf699d60fed3e775a9))
* Properly set nosurf context even when ignored ([0dcb774](https://github.com/ory/kratos/commit/0dcb774157bcbfd41a5d9df3914c31162226da75))
* Update cypress ([ba8b172](https://github.com/ory/kratos/commit/ba8b1729477233f79d099e5d7b397430ac1c6ace))
* Use correct regex for version replacement ([ce870ab](https://github.com/ory/kratos/commit/ce870ababdf089344a9428d3a405e18504a3c906)), closes [#787](https://github.com/ory/kratos/issues/787)

**Code Generation:**

* Pin v0.5.3-alpha.1 release commit ([64dc91a](https://github.com/ory/kratos/commit/64dc91af54cdf3eba158a50690240cdc8f7cb43b))

**Documentation:**

* Fix docosaurus admonitions ([#788](https://github.com/ory/kratos/issues/788)) ([281a7c9](https://github.com/ory/kratos/commit/281a7c9289570d4bee33447655281b610cbe7e52))
* Pin download script version ([e4137a6](https://github.com/ory/kratos/commit/e4137a6a41d68b1480af2075bda8c5f46c42cd22))
* Remove trailing garbage from quickstart ([#787](https://github.com/ory/kratos/issues/787)) ([7e70924](https://github.com/ory/kratos/commit/7e709242ada28b7781c6ace272f60f9d1b9d5b2f))

**Features:**

* Improve makefile install process and update deps ([d1eb37f](https://github.com/ory/kratos/commit/d1eb37f5d9d0f16e7864b5f8f08a44ba80853fa5))

**Tests:**

* Add e2e tests for mobile ([d481d51](https://github.com/ory/kratos/commit/d481d51f5f4de96cbbc7c347f5dbff381b44462d))
* Add option to disable csrf protection in apis ([a0077f1](https://github.com/ory/kratos/commit/a0077f12adf94ff428b502b69bbb0eaafd05be66))
* Bump wait time ([7a719e1](https://github.com/ory/kratos/commit/7a719e17c5641f4df47314f6f0ac2cf73dddc8bb))
* Install expo-cli globally ([db21cfa](https://github.com/ory/kratos/commit/db21cfa1c589a2dab829a4c8eaf1db15d14d965e))
* Install expo-cli in cci config with sudo ([d255f46](https://github.com/ory/kratos/commit/d255f462402f2d2c2278dcba1a139d0064343b22))
* Log wait-on output ([62b5ba9](https://github.com/ory/kratos/commit/62b5ba92d56e9f6b98adb8fb9c4daff03be08f2e))
* Output web server address ([cb41ca7](https://github.com/ory/kratos/commit/cb41ca78367b1943d230fa9ac116fcf3cf69b1c1))
* Resolve csrf test issues in settings ([ef8ba7d](https://github.com/ory/kratos/commit/ef8ba7dc93d6ba84f22b7aa65d00797e33b520a3))
* Resolve test panic ([6f6461f](https://github.com/ory/kratos/commit/6f6461fe3690576015ded9146c065a1e5d950be1))
* Revert delay increase and improve install scripts ([1eafcaa](https://github.com/ory/kratos/commit/1eafcaa86be194e412b0470a759bff6afc6c21af))






**Bug Fixes:**

* Add debug quickstart yml ([#780](https://github.com/ory/kratos/issues/780)) ([16e6b4d](https://github.com/ory/kratos/commit/16e6b4d76d297182ea9a1f5dc6367570f02f7b42))
* Gracefully handle double slashes in URLs ([aeb9414](https://github.com/ory/kratos/commit/aeb941477910b5ab54429a6aab7a3e1e388c48c5)), closes [#779](https://github.com/ory/kratos/issues/779)
* Merge gobuffalo CGO fix ([fea2e77](https://github.com/ory/kratos/commit/fea2e77ca0f9b20185c7a7704854fdcf29b7ab33))
* Remove obsolete recovery_token and add link to schema ([acf6ac4](https://github.com/ory/kratos/commit/acf6ac4e11c755e56c7d40728088257de367f7ff))
* Return correct error in login csrf ([dd9cab0](https://github.com/ory/kratos/commit/dd9cab0e02400c88e89877f755f03c6179013123)), closes [#785](https://github.com/ory/kratos/issues/785)
* Use correct assert package ([76be5b0](https://github.com/ory/kratos/commit/76be5b0a5d94c251f5f07eee9f700ec11b341e2e))

**Code Generation:**

* Pin v0.5.2-alpha.1 release commit ([79fcd8a](https://github.com/ory/kratos/commit/79fcd8a6949886f847f7be0c9ba2aba7554ab204))

**Documentation:**

* Small improvements to discord oidc provider guide ([#783](https://github.com/ory/kratos/issues/783)) ([6a3c453](https://github.com/ory/kratos/commit/6a3c45330885eb95015fa7ee9b58a72c38132499))

**Tests:**

* Add tests for csrf behavior ([48993e2](https://github.com/ory/kratos/commit/48993e2c496fb8af7e7b9e2752ba7078a134a75a)), closes [#785](https://github.com/ory/kratos/issues/785)
* Mark link as enabled in e2e test ([c214b81](https://github.com/ory/kratos/commit/c214b81a7026b06aaca062b2aa77951d01b0e237))
* Resolve schema test regression ([bb7af1b](https://github.com/ory/kratos/commit/bb7af1b759d6c812755956ef872bcbd31b9c50be))


