# Server Administration Guide

Table of Contents

- [Keycloak features and concepts](#keycloak-features-and-concepts)
  
  - [Features](#features)
  - [Basic Keycloak operations](#basic-keycloak-operations)
  - [Core concepts and terms](#core-concepts-and-terms)
- [Creating the first administrator](#creating-first-admin_server_administration_guide)
  
  - [Creating the account on the local host](#creating-the-account-on-the-local-host)
  - [Creating the account remotely](#creating-the-account-remotely)
- [Configuring realms](#_configuring-realms)
  
  - [Using the Admin Console](#using-the-admin-console)
  - [The master realm](#the-master-realm)
  - [Creating a realm](#proc-creating-a-realm_server_administration_guide)
  - [Configuring SSL for a realm](#_ssl_modes)
  - [Configuring email for a realm](#_email)
    
    - [XOAUTH2 email configuration with third-party vendors](#xoauth2-email-configuration-with-third-party-vendors)
  - [Configuring themes](#_themes)
  - [Enabling internationalization](#enabling-internationalization)
    
    - [User locale selection](#_user_locale_selection)
  - [Controlling login options](#controlling-login-options)
    
    - [Enabling forgot password](#enabling-forgot-password)
    - [Enabling Remember Me](#enabling-remember-me)
    - [ACR to Level of Authentication (LoA) Mapping](#_mapping-acr-to-loa-realm)
    - [Update Email Workflow (UpdateEmail)](#_update-email-workflow)
  - [Configuring realm keys](#realm_keys)
    
    - [Rotating keys](#rotating-keys)
    - [Adding a generated key pair](#adding-a-generated-key-pair)
    - [Rotating keys by extracting a certificate](#rotating-keys-by-extracting-a-certificate)
    - [Adding an existing key pair and certificate](#adding-an-existing-key-pair-and-certificate)
    - [Loading keys from a Java Keystore](#loading-keys-from-a-java-keystore)
    - [Making keys passive](#making-keys-passive)
    - [Disabling keys](#disabling-keys)
    - [Compromised keys](#compromised-keys)
- [Using external storage](#_user-storage-federation)
  
  - [Adding a provider](#adding-a-provider)
  - [Dealing with provider failures](#dealing-with-provider-failures)
  - [Lightweight Directory Access Protocol (LDAP) and Active Directory](#_ldap)
    
    - [Configuring federated LDAP storage](#configuring-federated-ldap-storage)
    - [Storage mode](#storage-mode)
    - [Edit mode](#edit-mode)
    - [Other configuration options](#other-configuration-options)
    - [Connecting to LDAP over SSL](#connecting-to-ldap-over-ssl)
    - [Synchronizing LDAP users to Keycloak](#synchronizing-ldap-users-to-keycloak)
    - [LDAP mappers](#_ldap_mappers)
    - [Password hashing](#_ldap_password_hashing)
    - [Enabling password change after reset](#_ldap_password_policy)
    - [Configuring the connection pool](#_ldap_connection_pool)
    - [Troubleshooting](#_ldap_troubleshooting)
  - [SSSD and FreeIPA Identity Management integration](#_sssd)
    
    - [FreeIPA/IdM server](#freeipaidm-server)
    - [SSSD and D-Bus](#sssd-and-d-bus)
    - [Enabling the SSSD federation provider](#enabling-the-sssd-federation-provider)
    - [Configuring a federated SSSD store](#configuring-a-federated-sssd-store)
  - [Custom providers](#custom-providers)
- [Managing users](#assembly-managing-users_server_administration_guide)
  
  - [Creating users](#proc-creating-user_server_administration_guide)
  - [Managing user attributes](#user-profile)
    
    - [Understanding the Default Configuration](#understanding-the-default-configuration)
    - [Understanding the User Profile Contexts](#understanding-the-user-profile-contexts)
    - [Understanding Managed and Unmanaged Attributes](#_understanding-managed-and-unmanaged-attributes)
    - [Managing the User Profile](#managing-the-user-profile)
    - [Managing Attributes](#managing-attributes)
    - [Validating Attributes](#_validating-attributes)
    - [Defining UI Annotations](#_defining-ui-annotations)
    - [Managing Attribute Groups](#managing-attribute-groups)
    - [Using the JSON configuration](#_user-profile-json-configuration)
    - [Customizing How UIs are Rendered](#customizing-how-uis-are-rendered)
    - [Enabling Progressive Profiling](#enabling-progressive-profiling)
    - [Using Internationalized Messages](#_using-internationalized-messages)
  - [Defining user credentials](#ref-user-credentials_server_administration_guide)
    
    - [Setting a password for a user](#proc-setting-password-user_server_administration_guide)
    - [Requesting a user reset a password](#requesting-a-user-reset-a-password)
    - [Creating an OTP](#proc_creating-otp_server_administration_guide)
  - [Allowing users to self-register](#con-user-registration_server_administration_guide)
    
    - [Enabling user registration](#proc-enabling-user-registration_server_administration_guide)
    - [Registering as a new user](#proc-registering-new-user_server_administration_guide)
    - [Requiring user to agree to terms and conditions during registration](#proc-requiring-tac-agreement-at-registration_server_administration_guide)
  - [Defining actions required at login](#con-required-actions_server_administration_guide)
    
    - [Setting required actions for one user](#proc-setting-required-actions_server_administration_guide)
    - [Setting required actions for all users](#proc-setting-default-required-actions_server_administration_guide)
    - [Enabling terms and conditions as a required action](#proc-enabling-terms-conditions_server_administration_guide)
  - [Application initiated actions](#con-aia_server_administration_guide)
    
    - [Re-authentication during AIA](#con-aia-reauth_server_administration_guide)
    - [Parameterized AIA](#con-aia-parameterized_server_administration_guide)
    - [Available actions](#con-aia-available-actions_server_administration_guide)
  - [Searching for a user](#proc-searching-user_server_administration_guide)
    
    - [Default search](#default-search)
    - [Attribute search](#attribute-search)
  - [Deleting a user](#proc-deleting-user_server_administration_guide)
  - [Enabling account deletion by users](#proc-allow-user-to-delete-account_server_administration_guide)
    
    - [Enabling the Delete Account Capability](#enabling-the-delete-account-capability)
    - [Giving a user the **delete-account** role](#giving-a-user-the-delete-account-role)
    - [Deleting your account](#deleting-your-account)
  - [Impersonating a user](#con-user-impersonation_server_administration_guide)
  - [Enabling reCAPTCHA](#proc-enabling-recaptcha_server_administration_guide)
    
    - [Setting up Google reCAPTCHA](#procedure_recaptcha)
    - [Setting up Google reCAPTCHA Enterprise](#procedure_recaptcha_enterprise)
  - [Personal data collected by Keycloak](#ref-personal-data-collected_server_administration_guide)
- [Managing user sessions](#managing-user-sessions)
  
  - [Administering sessions](#administering-sessions)
    
    - [Signing out all active sessions](#signing-out-all-active-sessions)
    - [Viewing client sessions](#viewing-client-sessions)
    - [Viewing user sessions](#viewing-user-sessions)
  - [Revoking active sessions](#_revocation-policy)
  - [Session and token timeouts](#_timeouts)
  - [Offline access](#_offline-access)
  - [Transient sessions](#_transient-session)
- [Assigning permissions using roles and groups](#assigning-permissions-using-roles-and-groups)
  
  - [Creating a realm role](#proc-creating-realm-roles_server_administration_guide)
  - [Client roles](#con-client-roles_server_administration_guide)
  - [Converting a role to a composite role](#_composite-roles)
  - [Assigning role mappings](#proc-assigning-role-mappings_server_administration_guide)
  - [Using default roles](#_default_roles)
  - [Role scope mappings](#_role_scope_mappings)
  - [Groups](#proc-managing-groups_server_administration_guide)
    
    - [Groups compared to roles](#con-comparing-groups-roles_server_administration_guide)
    - [Using default groups](#proc-specifying-default-groups_server_administration_guide)
- [Configuring authentication](#configuring-authentication_server_administration_guide)
  
  - [Password policies](#_password-policies)
    
    - [Password policy types](#password-policy-types)
  - [One Time Password (OTP) policies](#one-time-password-otp-policies)
    
    - [Time-based or counter-based one time passwords](#time-based-or-counter-based-one-time-passwords)
    - [TOTP configuration options](#totp-configuration-options)
    - [HOTP configuration options](#hotp-configuration-options)
  - [Authentication flows](#_authentication-flows)
    
    - [Built-in flows](#built-in-flows)
    - [Creating flows](#creating-flows)
    - [Creating a password-less browser login flow](#creating-a-password-less-browser-login-flow)
    - [Using Client Policies to Select an Authentication Flow](#_client-policy-auth-flow)
    - [Creating a browser login flow with step-up mechanism](#_step-up-flow)
    - [Step-up authentication for SAML](#_step-up-authentication-saml)
    - [Registration or Reset credentials requested by client](#_registration-rc-client-flows)
  - [User session limits](#_user_session_limits)
  - [Script Authenticator](#script-authenticator)
  - [Kerberos](#_kerberos)
    
    - [Setup of Kerberos server](#setup-of-kerberos-server)
    - [Setup and configuration of Keycloak server](#setup-and-configuration-of-keycloak-server)
    - [Setup and configuration of client machines](#setup-and-configuration-of-client-machines)
    - [Example setups](#example-setups)
    - [Credential delegation](#credential-delegation)
    - [Cross-realm trust](#cross-realm-trust)
    - [Troubleshooting](#troubleshooting)
  - [X.509 client certificate user authentication](#_x509)
    
    - [Features](#features-2)
    - [Adding X.509 client certificate authentication to browser flows](#_browser_flow)
    - [Configuring X.509 client certificate authentication](#_x509-config)
    - [Adding X.509 Client Certificate Authentication to a Direct Grant Flow](#adding-x-509-client-certificate-authentication-to-a-direct-grant-flow)
  - [W3C Web Authentication (WebAuthn)](#webauthn_server_administration_guide)
    
    - [Setup](#setup)
    - [Enable WebAuthn authentication in the default browser flow](#_webauthn-authenticator-setup)
    - [Authenticate with WebAuthn authenticator](#authenticate-with-webauthn-authenticator)
    - [Managing WebAuthn as an administrator](#managing-webauthn-as-an-administrator)
    - [Attestation statement verification](#attestation-statement-verification)
    - [Managing WebAuthn credentials as a user](#managing-webauthn-credentials-as-a-user)
    - [Registering WebAuthn credentials using AIA](#_webauthn_aia)
    - [Passwordless WebAuthn together with Two-Factor](#_webauthn_passwordless)
    - [LoginLess WebAuthn](#_webauthn_loginless)
  - [Passkeys](#passkeys_server_administration_guide)
    
    - [Passkey Authentication with Conditional UI or autofill](#_passkeys-conditional-ui)
    - [Passkeys Authentication with Modal UI](#passkeys-authentication-with-modal-ui)
    - [Setup](#setup-4)
  - [Recovery Codes](#_recovery-codes)
    
    - [Check Recovery Codes required action is enabled](#check-recovery-codes-required-action-is-enabled)
    - [Configure the Recovery Codes required action](#configure-the-recovery-codes-required-action)
    - [Adding Recovery Codes to the browser flow](#adding-recovery-codes-to-the-browser-flow)
    - [Creating the Recovery Codes credential](#creating-the-recovery-codes-credential)
  - [Conditions in conditional flows](#conditions-in-conditional-flows)
    
    - [Available conditions](#available-conditions)
    - [Explicitly deny/allow access in conditional flows](#explicitly-denyallow-access-in-conditional-flows)
    - [2FA conditional workflow examples](#twofa-conditional-workflow-examples)
  - [Authentication sessions](#_authentication-sessions)
    
    - [Authentication in more browser tabs](#authentication-in-more-browser-tabs)
- [Integrating identity providers](#_identity_broker)
  
  - [Brokering overview](#_identity_broker_overview)
  - [Default Identity Provider](#default_identity_provider)
  - [General configuration](#_general-idp-config)
  - [Social Identity Providers](#social-identity-providers)
    
    - [Bitbucket](#bitbucket)
    - [Facebook](#_facebook)
    - [GitHub](#_github)
    - [GitLab](#gitlab)
    - [Google](#_google)
    - [Instagram](#instagram)
    - [LinkedIn](#_linkedin)
    - [Microsoft](#_microsoft)
    - [OpenShift 4](#openshift-4)
    - [PayPal](#paypal)
    - [Stack Overflow](#_stackoverflow)
    - [Twitter](#_twitter)
  - [OpenID Connect v1.0 identity providers](#_identity_broker_oidc)
  - [OAuth v2 identity providers](#_identity_broker_oauth)
  - [SAML v2.0 Identity Providers](#saml-v2-0-identity-providers)
    
    - [Requesting specific AuthnContexts](#_identity_broker_saml_requested_authncontext)
    - [SP Descriptor](#_identity_broker_saml_sp_descriptor)
    - [Send subject in SAML requests](#_identity_broker_saml_login_hint)
  - [SPIFFE identity providers](#_identity_broker_spiffe)
  - [Kubernetes identity providers](#_identity_broker_kubernetes)
  - [Client-suggested Identity Provider](#_client_suggested_idp)
  - [Mapping claims and assertions](#_mappers)
  - [Available user session data](#available-user-session-data)
  - [First login flow](#_identity_broker_first_login)
    
    - [Default first login flow authenticators](#default-first-login-flow-authenticators)
    - [Automatically link existing first login flow](#automatically-link-existing-first-login-flow)
    - [Disabling automatic user creation](#_disabling_automatic_user_creation)
    - [Detect existing user first login flow](#_detect_existing_user_first_login_flow)
    - [Override existing broker link](#_override_existing_broker_link)
  - [Post login flow](#_identity_broker_post_login_flow)
    
    - [Post login flow examples](#post-login-flow-examples)
    - [Requesting additional authentication steps for the dedicated clients](#requesting-additional-authentication-steps-for-the-dedicated-clients)
  - [Retrieving external IDP tokens](#retrieving-external-idp-tokens)
  - [Identity broker logout](#identity-broker-logout)
- [SSO protocols](#sso-protocols)
  
  - [OpenID Connect](#con-oidc_server_administration_guide)
    
    - [OIDC auth flows](#con-oidc-auth-flows_server_administration_guide)
    - [OIDC Logout](#_oidc-logout)
    - [Keycloak server OIDC URI endpoints](#con-server-oidc-uri-endpoints_server_administration_guide)
  - [SAML](#_saml)
    
    - [SAML bindings](#con-saml-bindings_server_administration_guide)
    - [Keycloak Server SAML URI Endpoints](#keycloak-server-saml-uri-endpoints)
  - [OpenID Connect compared to SAML](#ref-saml-vs-oidc_server_administration_guide)
  - [Distribution Registry v2 authentication](#_docker)
    
    - [Docker authentication flow](#docker-authentication-flow)
    - [Keycloak Distribution registry v2 Authentication Server URI Endpoints](#keycloak-distribution-registry-v2-authentication-server-uri-endpoints)
- [Managing access to realm resources](#_admin_permissions)
  
  - [Master realm access control](#_master_realm_access_control)
    
    - [Global roles](#global-roles)
    - [Realm specific roles](#realm-specific-roles)
  - [Dedicated realm admin consoles](#_per_realm_admin_permissions)
  - [Delegating realm administration using permissions](#_fine_grained_permissions)
    
    - [Understanding the different types of realm administrators](#_understanding_different_types_realm_admins_)
    - [Understanding the Realm Resource Types](#understanding-the-realm-resource-types)
    - [Understanding the scopes of access](#understanding-the-scopes-of-access)
    - [Enabling admin permissions to a realm](#enabling-admin-permissions-to-a-realm)
    - [Managing Permissions](#_managing-permissions)
    - [Managing Policies](#managing-policies)
    - [Evaluating Permissions](#_evaluating-permissions)
    - [Accessing a Realm administration console as a Realm Administrator](#_realm_access_control)
    - [Understanding some common use cases](#understanding-some-common-use-cases)
    - [Performance considerations](#performance-considerations)
  - [Fine grained admin permissions V1](#fine-grained-admin-permissions-v1)
    
    - [Managing one specific client](#managing-one-specific-client)
    - [Restrict user role mapping](#restrict-user-role-mapping)
    - [Full list of permissions](#full-list-of-permissions)
- [Managing organizations](#_managing_organizations)
  
  - [Enabling organizations in Keycloak](#_enabling_organization_)
  - [Managing an organization](#managing-an-organization)
    
    - [Creating an organization](#creating-an-organization)
    - [Understanding organization domains](#understanding-organization-domains)
    - [Disabling an organization](#disabling-an-organization)
    - [Deleting an organization](#deleting-an-organization)
  - [Managing attributes](#_managing_attributes_)
  - [Managing members](#_managing_members_)
    
    - [Managed and unmanaged members](#_managed_unmanaged_members_)
    - [Adding an existing realm user as a member](#adding-an-existing-realm-user-as-a-member)
    - [Inviting users](#inviting-users)
    - [Managing invitations](#managing-invitations)
    - [Onboarding members using an Identity Provider](#_onboard_member_identity_provider_)
    - [Removing a member](#removing-a-member)
    - [Viewing organization group memberships](#viewing-organization-group-memberships)
    - [Support for federated members](#support-for-federated-members)
  - [Managing groups](#_managing_groups_)
    
    - [Creating groups](#creating-groups)
    - [Adding members to groups](#adding-members-to-groups)
    - [Understanding group paths](#understanding-group-paths)
    - [Mapping groups to tokens](#mapping-groups-to-tokens)
    - [Managing group attributes](#managing-group-attributes)
    - [Important distinctions](#important-distinctions)
    - [Deleting groups](#deleting-groups)
  - [Managing identity providers](#_managing_identity_provider_)
    
    - [Linking an identity provider to an organization](#linking-an-identity-provider-to-an-organization)
    - [Editing a linked identity provider](#editing-a-linked-identity-provider)
    - [Unlinking an identity provider from an organization](#unlinking-an-identity-provider-from-an-organization)
  - [Authenticating members](#authenticating-members_server_administration_guide)
    
    - [Understanding the identity-first login](#understanding-the-identity-first-login)
    - [Configuring existing authentication flows](#configuring-existing-authentication-flows)
    - [Configuring how users authenticate](#configuring-how-users-authenticate)
  - [Mapping organization claims](#_mapping_organization_claims_)
- [Managing workflows](#_managing_workflows)
  
  - [Understanding workflows](#_understanding_workflows_)
  - [Understanding the workflow definition](#_understanding_workflow_definition_)
  - [Understanding the workflow expression language](#_workflow_expression_language_)
  - [Managing workflows](#_managing_workflows_)
    
    - [Managing workflows through the Admin Console](#managing-workflows-through-the-admin-console)
  - [Triggering workflows on events](#_workflow_events_)
    
    - [Event functions](#_workflow_event_functions_)
  - [Scheduling workflows](#_scheduling_workflows_)
  - [Defining conditions](#_workflow_conditions_)
    
    - [User functions](#_workflow_user_functions_)
  - [Defining steps](#_workflow_steps_)
    
    - [User steps](#_workflow_user_steps_)
    - [Client steps](#_workflow_client_steps_)
    - [Understanding immediate steps](#_workflow_immediate_steps_)
    - [Understanding scheduled steps](#understanding-scheduled-steps)
  - [Understanding the workflows engine](#_understanding_workflows_engine_)
    
    - [Configuring the scheduled steps execution interval](#configuring-the-scheduled-steps-execution-interval)
    - [Configuring the task execution timeout](#configuring-the-task-execution-timeout)
    - [Performance considerations](#performance-considerations-2)
  - [Handling failures](#_handling_failures_)
  - [Troubleshooting workflows](#_troubleshooting_workflows_)
    
    - [Enabling workflow debug logging](#enabling-workflow-debug-logging)
    - [What to look for in the logs](#what-to-look-for-in-the-logs)
    - [Useful log categories](#useful-log-categories)
  - [Understanding common use cases](#_understanding_common_use_cases_)
    
    - [User Onboarding](#user-onboarding)
    - [User Offboarding](#user-offboarding)
    - [Tracking user activity and taking actions on inactivity](#tracking-user-activity-and-taking-actions-on-inactivity)
- [Managing OpenID Connect and SAML Clients](#assembly-managing-clients_server_administration_guide)
  
  - [Managing OpenID Connect clients](#_oidc_clients)
    
    - [Creating an OpenID Connect client](#proc-creating-oidc-client_server_administration_guide)
    - [Basic configuration](#con-basic-settings_server_administration_guide)
    - [Advanced configuration](#con-advanced-settings_server_administration_guide)
    - [Confidential client credentials](#_client-credentials)
    - [DPoP](#_dpop-bound-tokens)
    - [Client Secret Rotation](#_secret_rotation)
    - [Creating an OIDC Client Secret Rotation Policy](#_proc-secret-rotation)
    - [Using a service account](#_service_accounts)
    - [Role mappings in the token](#_oidc_token_role_mappings)
    - [Audience support](#audience-support)
  - [Creating a SAML client](#_client-saml-configuration)
    
    - [Settings tab](#settings-tab)
    - [Keys tab](#keys-tab)
    - [Advanced tab](#advanced-tab-2)
    - [IDP Initiated login](#idp-initiated-login)
    - [Using an entity descriptor to create a client](#proc-using-an-entity-descriptors_server_administration_guide)
  - [Client links](#con-client-links_server_administration_guide)
  - [OIDC token and SAML assertion mappings](#_protocol-mappers)
    
    - [Priority order](#_protocol-mappers_priority)
    - [OIDC user session note mappers](#_protocol-mappers_oidc-user-session-note-mappers)
    - [Script mapper](#script-mapper)
    - [Pairwise subject identifier mapper](#pairwise-subject-identifier-mapper)
    - [Using lightweight access token](#_using_lightweight_access_token)
  - [Generating client adapter config](#_client_installation)
  - [Client scopes](#_client_scopes)
    
    - [Protocol](#_client_scopes_protocol)
    - [Consent related settings](#consent-related-settings)
    - [Include in token scope](#include-in-token-scope)
    - [Link client scope with the client](#_client_scopes_linking)
    - [Evaluating Client Scopes](#_client_scopes_evaluate)
    - [Client scopes permissions](#client-scopes-permissions)
    - [Realm default client scopes](#realm-default-client-scopes)
    - [Downscoping](#_downscoping)
    - [Scopes explained](#scopes-explained)
  - [Client Policies](#_client_policies)
    
    - [Use-cases](#use-cases)
    - [Protocol](#protocol)
    - [Architecture](#architecture)
    - [Configuration](#configuration)
    - [Backward Compatibility](#backward-compatibility)
    - [Client Secret Rotation Example](#client-secret-rotation-example)
    - [Securing Client URIs](#securing-client-uris)
- [Configuring Keycloak as a Verifiable Credential Issuer](#_oid4vci)
  
  - [Introduction](#introduction)
  - [What are Verifiable Credentials (VCs)?](#what-are-verifiable-credentials-vcs)
  - [What is OID4VCI?](#what-is-oid4vci)
  - [Scope of This Chapter](#scope-of-this-chapter)
  - [Prerequisites](#prerequisites)
  - [Keycloak Instance](#keycloak-instance)
  - [Configuring Credential Issuance in Keycloak](#configuring-credential-issuance-in-keycloak)
  - [Authentication](#authentication)
  - [Configuration Steps](#configuration-steps)
  - [Creating a Realm](#creating-a-realm)
  - [Creating a User Account](#creating-a-user-account)
  - [Key Management Configuration](#key-management-configuration)
    
    - [Configuring Key Providers](#configuring-key-providers)
  - [Registering Realm Attributes](#registering-realm-attributes)
    
    - [Define Realm Attributes](#define-realm-attributes)
    - [Attribute Breakdown](#attribute-breakdown)
    - [Import Realm Attributes](#import-realm-attributes)
    - [Time-claim correlation mitigation](#time-claim-correlation-mitigation)
  - [Create Client Scopes with Mappers](#create-client-scopes-with-mappers)
    
    - [Define a Client Scope with a Mapper](#define-a-client-scope-with-a-mapper)
    - [Attribute Breakdown - ClientScope](#client-scope-attribute-breakdown)
    - [Attribute Breakdown - ProtocolMappers](#attribute-breakdown-protocolmappers)
    - [Import the Client Scope](#import-the-client-scope)
  - [Create the Client](#create-the-client)
  - [Verify the Configuration](#verify-the-configuration)
  - [Conclusion](#conclusion)
- [Using a vault to obtain secrets](#_vault-administration)
  
  - [Key resolvers](#_vault-key-resolvers)
- [Configuring auditing to track events](#configuring-auditing-to-track-events)
  
  - [Auditing user events](#auditing-user-events)
    
    - [Event types](#event-types)
    - [Event listener](#event-listener)
  - [Auditing admin events](#auditing-admin-events)
- [Mitigating security threats](#mitigating_security_threats)
  
  - [Host](#host)
  - [Admin endpoints and Admin Console](#admin-endpoints-and-admin-console)
  - [Brute force attacks](#password-guess-brute-force-attacks)
    
    - [Lockout permanently](#lockout-permanently)
    - [Lockout temporarily](#lockout-temporarily)
    - [Lockout permanently after temporary lockout](#lockout-permanently-after-temporary-lockout)
    - [Secondary Authentication Failures Lockout](#secondary-authentication-failures-lockout)
    - [Downside of Keycloak brute force detection](#downside-of-keycloak-brute-force-detection)
  - [Password policies](#password-policies)
  - [Read-only user attributes](#read_only_user_attributes)
  - [Validate user attributes](#validate_user_attributes)
  - [Clickjacking](#clickjacking)
  - [SSL/HTTPS requirement](#sslhttps-requirement)
  - [CSRF attacks](#csrf-attacks)
  - [Unspecific redirect URIs](#unspecific-redirect-uris_server_administration_guide)
  - [FAPI compliance](#fapi-compliance)
  - [OAuth 2.1 compliance](#oauth-2-1-compliance)
  - [Compromised access and refresh tokens](#compromised-access-and-refresh-tokens)
  - [Compromised authorization code](#compromised-authorization-code)
  - [Open redirectors](#open-redirectors)
  - [Mitigating Server-Side Request Forgery (SSRF)](#ssrf)
    
    - [Secure Client URIs Pattern Executor](#secure-client-uris-pattern-executor)
    - [Anti-SSRF Pattern Examples](#anti-ssrf-pattern-examples)
  - [Password database compromised](#password-database-compromised)
  - [Limiting scope](#limiting-scope)
    
    - [Scope availability](#scope-availability)
    - [Scope visibility](#scope-visibility)
  - [Limit token audience](#limit-token-audience)
  - [Limit Authentication Sessions](#_limit-authentication-sessions)
  - [SQL injection attacks](#sql-injection-attacks)
- [Account Console](#_account-service)
  
  - [Accessing the Account Console](#accessing-the-account-console)
  - [Configuring ways to sign in](#configuring-ways-to-sign-in)
    
    - [Two-factor authentication with OTP](#two-factor-authentication-with-otp)
    - [Two-factor authentication with WebAuthn](#two-factor-authentication-with-webauthn)
    - [Passwordless authentication with WebAuthn](#passwordless-authentication-with-webauthn)
  - [Viewing device activity](#viewing-device-activity)
  - [Adding an identity provider account](#adding-an-identity-provider-account)
  - [Accessing other applications](#accessing-other-applications)
  - [Viewing group memberships](#viewing-group-memberships)
- [Admin CLI](#admin-cli)
  
  - [Installing the Admin CLI](#installing-the-admin-cli)
  - [Using the Admin CLI](#using-the-admin-cli)
  - [Sensitive Options](#sensitive-options)
  - [Authenticating](#authenticating)
  - [Working with alternative configurations](#_working_with_alternative_configurations)
  - [Basic operations and resource URIs](#basic-operations-and-resource-uris)
  - [Realm operations](#realm-operations)
  - [Role operations](#role-operations)
  - [Client operations](#client-operations)
  - [User operations](#user-operations)
  - [Group operations](#_group_operations)
  - [Identity provider operations](#identity-provider-operations)
  - [Storage provider operations](#storage-provider-operations)
  - [Adding mappers](#adding-mappers)
  - [Authentication operations](#authentication-operations)

**Server Administration**

- [Getting Started](https://www.keycloak.org/guides#getting-started)
- [Securing applications](https://www.keycloak.org/guides#securing-apps)
- [Server Developer](https://www.keycloak.org/docs/26.6.3/server_development/)
- [Authorization Services](https://www.keycloak.org/docs/26.6.3/authorization_services/)
- [Upgrading](https://www.keycloak.org/docs/26.6.3/upgrading/)
- [Release Notes](https://www.keycloak.org/docs/26.6.3/release_notes/)

Version **26.6.3**

## [](#keycloak-features-and-concepts)Keycloak features and concepts

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/overview.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Foverview.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Foverview.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak is a single sign on solution for web apps and RESTful web services. The goal of Keycloak is to make security simple so that it is easy for application developers to secure the apps and services they have deployed in their organization. Security features that developers normally have to write for themselves are provided out of the box and are easily tailorable to the individual requirements of your organization. Keycloak provides customizable user interfaces for login, registration, administration, and account management. You can also use Keycloak as an integration platform to hook it into existing LDAP and Active Directory servers. You can also delegate authentication to third party identity providers like Facebook and Google.

### [](#features)Features

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/overview/features.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Foverview%2Ffeatures.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Foverview%2Ffeatures.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak provides the following features:

- Single Sign-On and Single Sign-Out for browser applications.
- OpenID Connect support.
- OAuth 2.0 support.
- SAML support.
- Identity Brokering - Authenticate with external OpenID Connect or SAML Identity Providers.
- Social Login - Enable login with Google, GitHub, Facebook, Twitter, and other social networks.
- User Federation - Sync users from LDAP and Active Directory servers.
- Kerberos bridge - Automatically authenticate users that are logged-in to a Kerberos server.
- Admin Console for central management of users, roles, role mappings, clients and configuration.
- Account Console that allows users to centrally manage their account.
- Theme support - Customize all user facing pages to integrate with your applications and branding.
- Flexible Authentication - Authenticate user with various mechanisms including passkey, password or X.509 certificates. Step-up authentication support
- Two-factor Authentication - Support for passkey, recovery codes and TOTP/HOTP via Google Authenticator or FreeOTP.
- Login flows - optional user self-registration, recover password, verify email, require password update, etc.
- Session management - Admins and users themselves can view and manage user sessions.
- Token mappers - Map user attributes, roles, etc. how you want into tokens and statements.
- Not-before revocation policies per realm, application and user.
- CORS support - Client adapters have built-in support for CORS.
- Service Provider Interfaces (SPI) - A number of SPIs to enable customizing various aspects of the server. Authentication flows, user federation providers, protocol mappers and many more.
- Supports any platform/language that has an OpenID Connect Relying Party library or SAML 2.0 Service Provider library.

### [](#basic-keycloak-operations)Basic Keycloak operations

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/overview/how.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Foverview%2Fhow.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Foverview%2Fhow.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak is a separate server that you manage on your network. Applications are configured to point to and be secured by this server. Keycloak uses open protocol standards like [OpenID Connect](https://openid.net/developers/how-connect-works/) or [SAML 2.0](https://saml.xml.org/saml-specifications) to secure your applications. Browser applications redirect a user’s browser from the application to the Keycloak authentication server where they enter their credentials. This redirection is important because users are completely isolated from applications and applications never see a user’s credentials. Applications instead are given an identity token or assertion that is cryptographically signed. These tokens can have identity information like username, address, email, and other profile data. They can also hold permission data so that applications can make authorization decisions. These tokens can also be used to make secure invocations on REST-based services.

### [](#core-concepts-and-terms)Core concepts and terms

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/overview/concepts.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Foverview%2Fconcepts.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Foverview%2Fconcepts.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Consider these core concepts and terms before attempting to use Keycloak to secure your web applications and REST services.

users

Users are entities that are able to log into your system. They can have attributes associated with themselves like email, username, address, phone number, and birthday. They can be assigned group membership and have specific roles assigned to them.

authentication

The process of identifying and validating a user.

authorization

The process of granting access to a user.

credentials

Credentials are pieces of data that Keycloak uses to verify the identity of a user. Some examples are passwords, one-time-passwords, digital certificates, or even fingerprints.

roles

Roles identify a type or category of user. `Admin`, `user`, `manager`, and `employee` are all typical roles that may exist in an organization. Applications often assign access and permissions to specific roles rather than individual users as dealing with users can be too fine-grained and hard to manage.

user role mapping

A user role mapping defines a mapping between a role and a user. A user can be associated with zero or more roles. This role mapping information can be encapsulated into tokens and assertions so that applications can decide access permissions on various resources they manage.

composite roles

A composite role is a role that can be associated with other roles. For example a `superuser` composite role could be associated with the `sales-admin` and `order-entry-admin` roles. If a user is mapped to the `superuser` role they also inherit the `sales-admin` and `order-entry-admin` roles.

groups

Groups manage groups of users. Attributes can be defined for a group. You can map roles to a group as well. Users that become members of a group inherit the attributes and role mappings that group defines.

realms

A realm manages a set of users, credentials, roles, and groups. A user belongs to and logs into a realm. Realms are isolated from one another and can only manage and authenticate the users that they control.

clients

Clients are entities that can request Keycloak to authenticate a user. Most often, clients are applications and services that want to use Keycloak to secure themselves and provide a single sign-on solution. Clients can also be entities that just want to request identity information or an access token so that they can securely invoke other services on the network that are secured by Keycloak.

client adapters

Client adapters are plugins that you install into your application environment to be able to communicate and be secured by Keycloak. Keycloak has a number of adapters for different platforms that you can download. There are also third-party adapters you can get for environments that we don’t cover.

consent

Consent is when you as an admin want a user to give permission to a client before that client can participate in the authentication process. After a user provides their credentials, Keycloak will pop up a screen identifying the client requesting a login and what identity information is requested of the user. User can decide whether or not to grant the request.

client scopes

When a client is registered, you must define protocol mappers and role scope mappings for that client. It is often useful to store a client scope, to make creating new clients easier by sharing some common settings. This is also useful for requesting some claims or roles to be conditionally based on the value of `scope` parameter. Keycloak provides the concept of a client scope for this.

client role

Clients can define roles that are specific to them. This is basically a role namespace dedicated to the client.

identity token

A token that provides identity information about the user. Part of the OpenID Connect specification.

access token

A token that can be provided as part of an HTTP request that grants access to the service being invoked on. This is part of the OpenID Connect and OAuth 2.0 specification.

assertion

Information about a user. This usually pertains to an XML blob that is included in a SAML authentication response that provided identity metadata about an authenticated user.

service account

Each client has a built-in service account which allows it to obtain an access token.

direct grant

A way for a client to obtain an access token on behalf of a user via a REST invocation.

protocol mappers

For each client you can tailor what claims and assertions are stored in the OIDC token or SAML assertion. You do this per client by creating and configuring protocol mappers.

session

When a user logs in, a session is created to manage the login session. A session contains information like when the user logged in and what applications have participated within single sign-on during that session. Both admins and users can view session information.

user federation provider

Keycloak can store and manage users. Often, companies already have LDAP or Active Directory services that store user and credential information. You can point Keycloak to validate credentials from those external stores and pull in identity information.

identity provider

An identity provider (IDP) is a service that can authenticate a user. Keycloak is an IDP.

identity provider federation

Keycloak can be configured to delegate authentication to one or more IDPs. Social login via Facebook or Google is an example of identity provider federation. You can also hook Keycloak to delegate authentication to any other OpenID Connect or SAML 2.0 IDP.

identity provider mappers

When doing IDP federation you can map incoming tokens and assertions to user and session attributes. This helps you propagate identity information from the external IDP to your client requesting authentication.

required actions

Required actions are actions a user must perform during the authentication process. A user will not be able to complete the authentication process until these actions are complete. For example, an admin may schedule users to reset their passwords every month. An `update password` required action would be set for all these users.

authentication flows

Authentication flows are work flows a user must perform when interacting with certain aspects of the system. A login flow can define what credential types are required. A registration flow defines what profile information a user must enter and whether something like reCAPTCHA must be used to filter out bots. Credential reset flow defines what actions a user must do before they can reset their password.

events

Events are audit streams that admins can view and hook into.

themes

Every screen provided by Keycloak is backed by a theme. Themes define HTML templates and stylesheets which you can override as needed.

## [](#creating-first-admin_server_administration_guide)Creating the first administrator

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/assembly-creating-first-admin.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fassembly-creating-first-admin.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fassembly-creating-first-admin.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

After installing Keycloak, you need an administrator account that can act as a *super* admin with full permissions to manage Keycloak. With this account, you can log in to the Keycloak Admin Console where you create realms and users and register applications that are secured by Keycloak.

### [](#creating-the-account-on-the-local-host)Creating the account on the local host

If your server is accessible from `localhost`, perform these steps.

Procedure

1. In a web browser, go to the [http://localhost:8080](http://localhost:8080) URL.
2. Supply a username and password that you can recall.
   
   Welcome page
   
   ![Welcome page](./images/initial-welcome-page.png)

### [](#creating-the-account-remotely)Creating the account remotely

If you cannot access the server from a `localhost` address or just want to start Keycloak from the command line, use the `KC_BOOTSTRAP_ADMIN_USERNAME` and `KC_BOOTSTRAP_ADMIN_PASSWORD` environment variables to create an initial admin account.

For example:

```
export KC_BOOTSTRAP_ADMIN_USERNAME=<username>
export KC_BOOTSTRAP_ADMIN_PASSWORD=<password>

bin/kc.[sh|bat] start
```

## [](#_configuring-realms)Configuring realms

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/admin-console.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fadmin-console.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fadmin-console.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Once you have an administrative account for the Admin Console, you can configure realms. A realm is a space where you manage objects, including users, applications, roles, and groups. A user belongs to and logs into a realm. One Keycloak deployment can define, store, and manage as many realms as there is space for in the database.

### [](#using-the-admin-console)Using the Admin Console

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/realms/proc-using-admin-console.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Frealms%2Fproc-using-admin-console.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Frealms%2Fproc-using-admin-console.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You configure realms and perform most administrative tasks in the Keycloak Admin Console.

Prerequisites

To use the Admin Console, you need an administrator account.

- If no administrators exist, see [Creating the first administrator](#creating-first-admin_server_administration_guide).
- If other administrators exist, ask an administrator to provide an account with privileges to manage realms.

Procedure

1. Go to the URL for the Admin Console.
   
   For example, for localhost, use this URL: [http://localhost:8080/admin/](http://localhost:8080/admin/)
2. Enter the username and password you created on the Welcome Page or through environment variables as described in [Creating the initial admin user](https://www.keycloak.org/server/configuration#_creating_the_initial_admin_user).
   
   Login page
   
   ![Login page](./images/login-page.png)
   
   This action displays the Admin Console.
   
   Admin Console
   
   ![Admin Console](./images/admin-console.png)
3. Note the menus and other options that you can use:
   
   - Click the **Current realm** to see if other realms are available to be managed.
   - Click **Create realm** to create another realm that you can manage.
   - Click the top right list to view your account or log out.
4. Click **Realm settings** in the menu to see the fields and options for this realm.
   
   Click a question mark **?** icon to show the definition of a field such as **Frontend URL**.
   
   Realm settings
   
   ![Realm settings](./images/realm-settings.png)

Export files from the Admin Console are not suitable for backups or data transfer between servers. Only boot-time exports are suitable for backups or data transfer between servers.

### [](#the-master-realm)The master realm

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/realms/master.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Frealms%2Fmaster.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Frealms%2Fmaster.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

In the Admin Console, two types of realms exist:

- `Master realm` - This realm was created for you when you first started Keycloak. It contains the administrator account you created at the first login. Use the *master* realm only to create and manage the realms in your system.
- `Other realms` - These realms are created by the administrator in the master realm. In these realms, administrators manage the users in your organization and the applications they need. The applications are owned by the users.

Realms and applications

![Realms and applications](./images/master_realm.png)

Realms are isolated from one another and can only manage and authenticate the users that they control. Following this security model helps prevent accidental changes and follows the tradition of permitting user accounts access to only those privileges and powers necessary for the successful completion of their current task.

Additional resources

- See [Dedicated Realm Admin Consoles](#_per_realm_admin_permissions) if you want to disable the *master* realm and define administrator accounts within any new realm you create. Each realm has its own dedicated Admin Console that you can log into with local accounts.

### [](#proc-creating-a-realm_server_administration_guide)Creating a realm

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/realms/proc-creating-a-realm.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Frealms%2Fproc-creating-a-realm.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Frealms%2Fproc-creating-a-realm.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You create a realm to provide a management space where you can create users and give them permissions to use applications. At first login, you are typically in the *master* realm, the top-level realm from which you create other realms.

When deciding what realms you need, consider the kind of isolation you want to have for your users and applications. For example, you might create a realm for the employees of your company and a separate realm for your customers. Your employees would log into the employee realm and only be able to visit internal company applications. Customers would log into the customer realm and only be able to interact with customer-facing apps.

Procedure

1. In the Admin Console, click **Create Realm** next to **Current realm**.
2. Enter a name for the realm.
3. Click **Create**.
   
   Create realm
   
   ![Create realm](./images/create-realm.png)
   
   The current realm is now set to the realm you just created. You can switch between realms by clicking the realm name in the menu.

### [](#_ssl_modes)Configuring SSL for a realm

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/realms/ssl.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Frealms%2Fssl.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Frealms%2Fssl.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Each realm has an associated SSL Mode, which defines the SSL/HTTPS requirements for interacting with the realm. Browsers and applications that interact with the realm honor the SSL/HTTPS requirements defined by the SSL Mode or they cannot interact with the server.

Procedure

1. Click **Realm settings** in the menu.
2. Click the **General** tab.
   
   General tab
   
   ![General Tab](./images/general-tab.png)
3. Set **Require SSL** to one of the following SSL modes:
   
   - **External requests** Users can interact with Keycloak without SSL so long as they stick to private IPv4 addresses such as `localhost`, `127.0.0.1`, `10.x.x.x`, `192.168.x.x`, `172.16.x.x` or IPv6 link-local and unique-local addresses. If you try to access Keycloak without SSL from a non-private IP address, you will get an error.
   - **None** Keycloak does not require SSL. This choice applies only in development when you are experimenting and do not plan to support this deployment.
   - **All requests** Keycloak requires SSL for all IP addresses.

### [](#_email)Configuring email for a realm

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/realms/email.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Frealms%2Femail.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Frealms%2Femail.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak sends emails to users to verify their email addresses, when they forget their passwords, or when an administrator needs to receive notifications about a server event. To enable Keycloak to send emails, you provide Keycloak with your SMTP server settings.

Procedure

1. Click **Realm settings** in the menu.
2. Click the **Email** tab.
   
   Email tab
   
   ![Email Tab](./images/email-tab.png)
3. Fill in the fields and toggle the switches as needed.

Template

From

**From** denotes the address used for the **From** SMTP-Header for the emails sent.

From display name

**From display name** allows to configure a user-friendly email address aliases (optional). If not set the plain **From** email address will be displayed in email clients.

Reply to

**Reply to** denotes the address used for the **Reply-To** SMTP-Header for the mails sent (optional). If not set the plain **From** email address will be used.

Reply to display name

**Reply to display name** allows to configure a user-friendly email address aliases (optional). If not set the plain **Reply To** email address will be displayed.

Envelope from

**Envelope from** denotes the [Bounce Address](https://en.wikipedia.org/wiki/Bounce_address) used for the **Return-Path** SMTP-Header for the mails sent (optional).

Connection & Authentication

Host

**Host** denotes the SMTP server hostname used for sending emails.

Port

**Port** denotes the SMTP server port.

Encryption

Tick one of these checkboxes to support sending emails for recovering usernames and passwords, especially if the SMTP server is on an external network. You will most likely need to change the **Port** to 465, the default port for SSL/TLS.

Authentication

Set this switch to **ON** if your SMTP server requires authentication.

Username

All authentication-mechanisms require a username.

Authentication Type

Choose the kind of authentication: 'password' or 'token'.

Password

Only needed when **Authentication Type** 'password' is selected. Supply the **Password**. The value of the **Password** field can refer a value from an external [vault](#_vault-administration).

Auth Token URL

Only needed when **Authentication Type** 'token' is selected. Supply the **Auth Token URL** that is used to fetch a token via client credentials grant.

Auth Token Scope

Only needed when **Authentication Type** 'token' is selected. Supply the **Auth Token Scope** that is used to fetch a token from the **Auth Token URL**.

Auth Token ClientId

Only needed when **Authentication Type** 'token' is selected. Supply the **Auth ClientId** that is used to fetch a token from the **Auth Token URL**.

Auth Token Client Secret

Only needed when **Authentication Type** 'token' is selected. Supply the **Auth Client Secret** that authenticates the client to fetch a token from the **Auth Token URL**. The value of the **Auth Client Secret** field can refer a value from an external [vault](#_vault-administration).

Allow UTF-8

Enable to UTF-8-encode email address when sending them to the server. This should only be enabled if the mail server supports UTF-8 via the SMTPUTF8 extension. If disabled, domain names containing non-ASCII characters will be encoded using punycode, and addresses containing non-ASCII characters in the local part of the address will return an error.

If the realm is configured to send emails (this SMTP configuration is setup) and **Allow UTF-8** option is disabled, the built-in [user profile](#user-profile) email validator checks the local part of the address contains only ASCII characters. This way, Keycloak prevents user emails that cannot be notified.

#### [](#xoauth2-email-configuration-with-third-party-vendors)XOAUTH2 email configuration with third-party vendors

The following section contains some hints on how to configure Keycloak email settings to use XOAUTH2 based authentication with some known third-party software SMTP servers.

This section has been contributed by the Keycloak community. As the Keycloak core team does not have means to test third-party providers, it is provided as-is. If you find this documentation outdated or incomplete, please contribute to improve it.

##### [](#configuration-for-microsoft-azure-and-office365)Configuration for Microsoft Azure and Office365

Microsoft Azure allows 'Client Credentials Grant' using a client secret to gather an access token. Microsoft Office365 supports SMTP with XOAUTH2 to authenticate with the gathered token.

Links to relevant Microsoft documentation:

- [Usage of role base access control for applications in exchange online](https://learn.microsoft.com/en-us/exchange/permissions-exo/application-rbac)
- Settings in [Authenticate an IMAP, POP or SMTP connection using OAuth](https://learn.microsoft.com/en-us/exchange/client-developer/legacy-protocols/how-to-authenticate-an-imap-pop-smtp-application-by-using-oauth)

The following method for setting up Keycloak to send email with Azure and Office365 has been verified by a test. There might be other variants to achieve the same depending on your environment.

From

`<some>@<domain>`

Host

`smtp.office365.com`

Port

`587`

Encryption

Check Start TLS

Username

`<some>@<domain>` (might be the same of a different value than the sender value)

Auth Token Url

`https://login.microsoftonline.com/<TenantID>/oauth2/v2.0/token`

Replace TenantID with the id of your Microsoft tenant, usually a UUID, in Azure or just copy the token url from the list of endpoints displayed in the Azure Console.

Auth Token Scope

`https://outlook.office.com/.default`

Auth Token ClientId

`<ApplicationId>`

Replace ApplicationId with the id of your application in Azure, usually a UUID.

Auth Token ClientSecret

`<Secret configured>`

##### [](#configuration-for-google-mail)Configuration for Google Mail

This feature is not yet supported by Keycloak, because Google does not allow client-secrets for the Client Credentials Grant.

##### [](#configuration-for-aws)Configuration for AWS

XOAUTH2 is not supported by the AWS-SMTP service. The AWS-service requires the use of a password.

### [](#_themes)Configuring themes

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/realms/themes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Frealms%2Fthemes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Frealms%2Fthemes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

For a given realm, you can change the appearance of any UI in Keycloak by using themes.

Procedure

1. Click **Realm settings** in the menu.
2. Click the **Themes** tab.
   
   Themes tab
   
   ![Themes tab](./images/themes-tab.png)
3. Pick the theme you want for each UI category and click **Save**.
   
   Login theme
   
   Username password entry, OTP entry, new user registration, and other similar screens related to login.
   
   Account theme
   
   The console used by the user to manage his or her account.
   
   Admin console theme
   
   The skin of the Keycloak Admin Console.
   
   Email theme
   
   Whenever Keycloak has to send out an email, it uses templates defined in this theme to craft the email.

Additional resources

- For details on creating or modifying themes, see [Deploying Themes](https://www.keycloak.org/ui-customization/themes#_deploying_themes) in the UI Customization Guide.

### [](#enabling-internationalization)Enabling internationalization

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/realms/proc-configuring-internationalization.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Frealms%2Fproc-configuring-internationalization.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Frealms%2Fproc-configuring-internationalization.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Every UI screen is internationalized in Keycloak. The default language is English, but you can choose which locales you want to support and what the default locale will be.

Procedure

1. Click **Realm Settings** in the menu.
2. Click the **Localization** tab.
3. Enable **Internationalization**.
4. Select the languages you will support.
   
   Localization tab
   
   ![Localization tab](./images/localization.png)
   
   The next time a user logs in, that user can choose a language on the login page to use for the login screens, Account Console, and Admin Console.

Additional resources

- The [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) explains how you can offer additional languages. All internationalized texts which are provided by the theme can be overwritten by realm-specific texts on the **Localization** tab.

#### [](#_user_locale_selection)User locale selection

A locale selector provider suggests the best locale on the information available. However, it is often unknown who the user is. For this reason, the previously authenticated user’s locale is remembered in a persisted cookie.

The logic for selecting the locale uses the first of the following that is available:

- User selected - when the user has selected a locale using the drop-down locale selector
- User profile - when there is an authenticated user and the user has a preferred locale set
- Client selected - passed by the client using for example ui\_locales parameter
- Cookie - last locale selected on the browser
- Accepted language - locale from **Accept-Language** header
- Realm default
- If none of the above, fall back to English

When a user is authenticated an action is triggered to update the locale in the persisted cookie mentioned earlier. If the user has actively switched the locale through the locale selector on the login pages the users locale is also updated at this point.

If you want to change the logic for selecting the locale, you have an option to create custom `LocaleSelectorProvider`. For details, please refer to the [Working with themes: Locale selector](https://www.keycloak.org/ui-customization/themes).

### [](#controlling-login-options)Controlling login options

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/login-settings.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Flogin-settings.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Flogin-settings.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak includes several built-in login page features.

#### [](#enabling-forgot-password)Enabling forgot password

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/login-settings/forgot-password.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Flogin-settings%2Fforgot-password.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Flogin-settings%2Fforgot-password.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

If you enable `Forgot password`, users can reset their login credentials if they forget their passwords or lose their OTP generator.

Procedure

1. Click **Realm settings** in the menu.
2. Click the **Login** tab.
   
   Login tab
   
   ![Login Tab](./images/login-tab.png)
3. Toggle **Forgot password** to **ON**.
   
   A `Forgot Password?` link displays in your login pages.
   
   Forgot password link
   
   ![Forgot Password Link](./images/forgot-password-link.png)
4. Specify `Host` and `From` in the **Email** tab in order for Keycloak to be able to send the reset email.
5. Click this link to bring users where they can enter their username or email address and receive an email with a link to reset their credentials.
   
   Forgot password page
   
   ![Forgot Password Page](./images/forgot-password-page.png)

The text sent in the email is configurable. See [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more information.

When users click the email link, Keycloak asks them to update their password, and if they have set up an OTP generator, Keycloak asks them to reconfigure the OTP generator. For security reasons, the flow forces federated users to login again after the reset credentials and keeps internal database users logged in if the same authentication session (same browser) is used. Depending on the security requirements of your organization, you can change the default behavior.

To change this behavior, perform these steps:

Procedure

1. Click **Authentication** in the menu.
2. Click the **Flows** tab.
3. Select the **Reset Credentials** flow.
   
   Reset credentials flow
   
   ![Reset Credentials Flow](./images/reset-credentials-flow.png)
   
   If you do not want to reset the OTP, set the `Reset - Conditional OTP` sub-flow requirement to **Disabled**.
   
   Send Reset Email Configuration
   
   ![Send Reset Email Configuration](./images/reset-credential-email-config.png)
   
   If you want to change default behavior for the force login option, click the **Send Reset Email** settings icon in the flow, define an **Alias**, and select the best **Force login after reset** option for you (`true`, always force re-authentication, `false`, keep the user logged in if the same browser was used, `only-federated`, default value that forces login again only for federated users).
4. Click **Authentication** in the menu.
5. Click the **Required actions** tab.
6. Ensure **Update Password** is enabled.
   
   Required Actions
   
   ![Required Actions](./images/reset-credentials-required-actions.png)

#### [](#enabling-remember-me)Enabling Remember Me

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/login-settings/remember-me.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Flogin-settings%2Fremember-me.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Flogin-settings%2Fremember-me.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

A logged-in user closing their browser destroys their session, and that user must log in again. You can set Keycloak to keep the user’s login session open if that user clicks the *Remember Me* checkbox upon login. This action turns the login cookie from a session-only cookie to a persistence cookie.

Procedure

1. Click **Realm settings** in the menu.
2. Click the **Login** tab.
3. Toggle the **Remember Me** switch to **On**.
   
   Login tab
   
   ![Login Tab Remember Me](./images/login-tab-remember-me.png)
   
   When you save this setting, a `remember me` checkbox displays on the realm’s login page.
   
   Remember Me
   
   ![Remember Me](./images/remember-me.png)

Disabling the "Remember me" option will invalidate all sessions created with the "Remember me" checkbox selected during login, requiring users to log in again. Any refresh tokens related to these sessions will also become invalid.

The sessions will not be invalidated immediately when the switch is disabled, but when a cookie or token associated with an invalid session is used, or asynchronously in the background. This means that disabling and then re-enabling the "Remember me" switch cannot be used to invalidate old sessions.

#### [](#_mapping-acr-to-loa-realm)ACR to Level of Authentication (LoA) Mapping

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/login-settings/acr-to-loa-mapping.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Flogin-settings%2Facr-to-loa-mapping.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Flogin-settings%2Facr-to-loa-mapping.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

In the general settings of a realm, you can define which `Authentication Context Class Reference (ACR)` value is mapped to which `Level of Authentication (LoA)`. The ACR can be any value, whereas the LoA must be numeric. The acr claim can be requested in the `claims` or `acr_values` parameter sent in the OIDC request and it is also included in the access token and ID token. The mapped number is used in the authentication flow conditions.

Mapping can be also specified at the client level in case that particular client needs to use different values than realm. However, a best practice is to stick to realm mappings.

ACR to LoA mapping

![ACR to LoA mapping](./images/realm-oidc-map-acr-to-loa.png)

For further details see [Step-up Authentication](#_step-up-flow) and [the official OIDC specification](https://openid.net/specs/openid-connect-core-1_0.html#acrSemantics).

If the feature for [Step-up authentication for SAML](#_step-up-authentication-saml) is enabled, the ACR to LoA mapping is a table with three values. The new URI column is the URI that will map the SAML authentication context class reference to the numeric LoA. This new column is necessary if you want to use step-up authentication with the SAML protocol.

ACR/URI to LoA mapping

![ACR/URI to LoA mapping](./images/realm-oidc-map-acr-uri-to-loa.png)

#### [](#_update-email-workflow)Update Email Workflow (UpdateEmail)

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/login-settings/update-email-workflow.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Flogin-settings%2Fupdate-email-workflow.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Flogin-settings%2Fupdate-email-workflow.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

With this workflow, users will have to use an `UPDATE_EMAIL` action to change their own email address.

This action provides a more secure and consistent flow to update user emails by requiring re-authentication and optionally requiring email verification before any update to their account.

Applications are able to send their users to the email update form by leveraging UPDATE\_EMAIL as an [AIA (Application Initiated Action)](#con-aia_server_administration_guide).

To enable Update Email capability for a realm, go to the `Authentication` menu in the administration console and click on `Required actions` tab. Switch the toggle for `UPDATE_EMAIL` required action to `enabled`.

##### [](#forcing-users-to-re-authenticate-before-updating-email)Forcing users to re-authenticate before updating email

When the `UPDATE_EMAIL` required action is enabled, users may be required to re-authenticate before being able to update their email if their last authentication is older than the configured duration. This is a security measure to prevent account takeover in case the user credentials are not known by the attacker but the user session is hijacked.

By default, the user will be asked to re-authenticate if the last authentication is older than 5 minutes (300 seconds). You can change this value by setting the `Maximum Age of Authentication` setting in the `UPDATE_EMAIL` required action configuration. By setting this value to `0`, the user will always be asked to re-authenticate before updating the email.

##### [](#verifying-emails)Verifying Emails

If the realm has email verification disabled, this action will allow to update the email without verification.

If the realm has email verification enabled, the action will send an email with a link to the new email address without changing the account email. Only after following the link and confirming the email, the email will be updated.

Under certain circumstances, you do not want to enable email verification at the realm level but only when users are updating their emails. For that, you can set the `Force Email Verification` setting on the `UPDATE_EMAIL` required action to force users to verify their emails even though email verification is eventually disabled at the realm level. By default, email verification is not enabled.

In case the user is updating the email during the authentication flow (e.g.: when running the `UPDATE_PROFILE` required action), the user will be forced to verify the email if any of the `Verify Email` or the `Force Email Verification` settings are enabled. In case the `Verify Email` is enabled at the realm level, the `VERIFY_EMAIL` required action will be automatically added to the user account. Otherwise, if only the `Force Email Verification` is enabled the `UPDATE_EMAIL` required action will be added instead.

If a user has `Email Verified` set, and both `Verify Email` and `Force Email Verification` are disabled, `Email Verified` resets after the user updates email.

##### [](#updating-the-user-email)Updating the user email

When the `UPDATE_EMAIL` required action is enabled, the user can update their emails by:

- Self-registering to a realm if this capability is enabled to realm
- Accessing the account console and clicking the `Update email` link when at the `Personal info` section
- Updating the profile during the authentication flow (e.g.: when running the `UPDATE_PROFILE` required action) if the email is not yet set. If an existing user does have an email set when updating the profile during the authentication flow, the email attribute will not be available.
- Administrators when updating the user account through the administration console

##### [](#pending-email-verification)Pending Email Verification

When a user initiates an email update that requires verification, the new email address is stored in a pending state until the user clicks the verification link. If the user tries to log in again before clicking the verification link, they will see a message informing them that a verification email was sent to the new address, with options to resend the email or enter a different email address.

Administrators can manage these pending verifications through the admin console. In the user details page, when a user has a pending email verification, a warning alert is displayed indicating the pending verification status. The alert shows which email address is awaiting confirmation and provides a link to cancel the verification process.

Clicking this link opens a confirmation dialog that allows administrators to remove the pending verification or cancel the action.

When confirmed, this action will:

- Remove the pending email verification attribute
- Invalidate the existing verification link
- Remove the `UPDATE_EMAIL` required action from the user

##### [](#update-email-and-user-profile)Update Email and User Profile

If the email attribute is set as required in the user profile configuration, the requirement is kept in the Update Email workflow, meaning a user won’t be able to clear his/her email in update email page. The opposite is true, if the email attribute is set as optional in the user profile configuration.

If the email attribute is set as read-only in the user profile configuration, the following behavior applies:

- The `Update email` link will not be displayed in the account console
- The `UPDATE_EMAIL` required action will be automatically skipped and removed from the user account
- In the update profile page, if the user’s email is initially empty, the email field will be hidden

##### [](#message-customization)Message Customization

All messages displayed in this workflow, including admin console messages, verification emails, and update email page messages, can be customized using the standard Keycloak message customization system.

### [](#realm_keys)Configuring realm keys

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/realms/keys.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Frealms%2Fkeys.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Frealms%2Fkeys.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The authentication protocols that are used by Keycloak require cryptographic signatures and sometimes encryption. Keycloak uses asymmetric key pairs, a private and public key, to accomplish this.

Keycloak has a single active key pair at a time, but can have several passive keys as well. The active key pair is used to create new signatures, while the passive key pair can be used to verify previous signatures. This makes it possible to regularly rotate the keys without any downtime or interruption to users.

When a realm is created, a key pair and a self-signed certificate is automatically generated.

Procedure

1. Click **Realm settings** in the menu.
2. Click **Keys**.
3. Select **Passive keys** from the filter dropdown to view passive keys.
4. Select **Disabled keys** from the filter dropdown to view disabled keys.

A key pair can have the status `Active`, but still not be selected as the currently active key pair for the realm. The selected active pair which is used for signatures is selected based on the first key provider sorted by priority that is able to provide an active key pair.

#### [](#rotating-keys)Rotating keys

We recommend that you regularly rotate keys. Start by creating new keys with a higher priority than the existing active keys. You can instead create new keys with the same priority and making the previous keys passive.

Once new keys are available, all new tokens and cookies will be signed with the new keys. When a user authenticates to an application, the SSO cookie is updated with the new signature. When OpenID Connect tokens are refreshed new tokens are signed with the new keys. Eventually, all cookies and tokens use the new keys and after a while the old keys can be removed.

The frequency of deleting old keys is a tradeoff between security and making sure all cookies and tokens are updated. Consider creating new keys every three to six months and deleting old keys one to two months after you create the new keys. If a user was inactive in the period between the new keys being added and the old keys being removed, that user will have to re-authenticate.

Rotating keys also applies to offline tokens. To make sure they are updated, the applications need to refresh the tokens before the old keys are removed.

#### [](#adding-a-generated-key-pair)Adding a generated key pair

Use this procedure to generate a key pair including a self-signed certificate.

Procedure

1. Select the realm in the Admin Console.
2. Click **Realm settings** in the menu.
3. Click the **Keys** tab.
4. Click the **Providers** tab.
5. Click **Add provider** and select **rsa-generated**.
6. Enter a number in the **Priority** field. This number determines if the new key pair becomes the active key pair. The highest number makes the key pair active.
7. Select a value for **AES Key size**.
8. Click **Save**.

Changing the priority for a provider will not cause the keys to be re-generated, but if you want to change the keysize you can edit the provider and new keys will be generated.

#### [](#rotating-keys-by-extracting-a-certificate)Rotating keys by extracting a certificate

You can rotate keys by extracting a certificate from an RSA generated key pair and using that certificate in a new keystore.

Prerequisites

- A generated key pair

Procedure

1. Select the realm in the Admin Console.
2. Click **Realm Settings**.
3. Click the **Keys** tab.
   
   A list of **Active** keys appears.
4. On a row with an RSA key, click **Certificate** under **Public Keys**.
   
   The certificate appears in text form.
5. Save the certificate to a file and enclose it in these lines.
   
   ```
   ----Begin Certificate----
   <Output>
   ----End Certificate----
   ```
6. Use the **keytool** command to convert the key file to PEM Format.
7. Remove the current RSA public key certificate from the keystore.
   
   ```
   keytool -delete -keystore <keystore>.jks -storepass <password> -alias <key>
   ```
8. Import the new certificate into the keystore
   
   ```
   keytool -importcert -file domain.crt -keystore <keystore>.jks -storepass <password>  -alias <key>
   ```
9. Rebuild the application.
   
   ```
   mvn clean install wildfly:deploy
   ```

#### [](#adding-an-existing-key-pair-and-certificate)Adding an existing key pair and certificate

To add a key pair and certificate obtained elsewhere select `Providers` and choose `rsa` from the dropdown. You can change the priority to make sure the new key pair becomes the active key pair.

Prerequisites

- A private key file. The file must be PEM formatted.

Procedure

1. Select the realm in the Admin Console.
2. Click **Realm settings**.
3. Click the **Keys** tab.
4. Click the **Providers** tab.
5. Click **Add provider** and select **rsa**.
6. Enter a number in the **Priority** field. This number determines if the new key pair becomes the active key pair.
7. Click **Browse…​** beside **Private RSA Key** to upload the private key file.
8. If you have a signed certificate for your private key, click **Browse…​** beside **X509 Certificate** to upload the certificate file. Keycloak automatically generates a self-signed certificate if you do not upload a certificate.
9. Click **Save**.

#### [](#loading-keys-from-a-java-keystore)Loading keys from a Java Keystore

To add a key pair and certificate stored in a Java Keystore file on the host select `Providers` and choose `java-keystore` from the dropdown. You can change the priority to make sure the new key pair becomes the active key pair.

For the associated certificate chain to be loaded it must be imported to the Java Keystore file with the same `Key Alias` used to load the key pair.

Procedure

01. Select the realm in the Admin Console.
02. Click **Realm settings** in the menu.
03. Click the **Keys** tab.
04. Click the **Providers** tab.
05. Click **Add provider** and select **java-keystore**.
06. Enter a number in the **Priority** field. This number determines if the new key pair becomes the active key pair.
07. Enter the desired **Algorithm**. Note that the algorithm should match the key type (for example `RS256` requires a RSA private key, `ES256` a EC private key or `AES` an AES secret key).
08. Enter a value for **Keystore**. Path to the keystore file.
09. Enter the **Keystore Password**. The option can refer a value from an external [vault](#_vault-administration).
10. Enter a value for **Keystore Type** (`JKS`, `PKCS12` or `BCFKS`).
11. Enter a value for the **Key Alias** to load from the keystore.
12. Enter the **Key Password**. The option can refer a value from an external [vault](#_vault-administration).
13. Enter a value for **Key Use** (`sig` for signing or `enc` for encryption). Note that the use should match the algorithm type (for example `RS256` is `sig` but `RSA-OAEP` is `enc`)
14. Click **Save**.

Not all the keystore types support all types of keys. For example, `JKS` in all modes and `PKCS12` in fips mode (`BCFIPS` provider) cannot store secret key entries.

#### [](#making-keys-passive)Making keys passive

Procedure

1. Select the realm in the Admin Console.
2. Click **Realm settings** in the menu.
3. Click the **Keys** tab.
4. Click the **Providers** tab.
5. Click the provider of the key you want to make passive.
6. Toggle **Active** to **Off**.
7. Click **Save**.

#### [](#disabling-keys)Disabling keys

Procedure

1. Select the realm in the Admin Console.
2. Click **Realm settings** in the menu.
3. Click the **Keys** tab.
4. Click the **Providers** tab.
5. Click the provider of the key you want to make passive.
6. Toggle **Enabled** to **Off**.
7. Click **Save**.

#### [](#compromised-keys)Compromised keys

Keycloak has the signing keys stored just locally and they are never shared with the client applications, users or other entities. However, if you think that your realm signing key was compromised, you should first generate new key pair as described above and then immediately remove the compromised key pair.

Alternatively, you can delete the provider from the `Providers` table.

Procedure

1. Click **Clients** in the menu.
2. Click **security-admin-console**.
3. Scroll down to the **Access settings** section.
4. Fill in the **Admin URL** field.
5. Click the **Advanced** tab.
6. Click **Set to now** in the **Revocation** section.
7. Click **Push**.

Pushing the not-before policy ensures that client applications do not accept the existing tokens signed by the compromised key. The client application is forced to download new key pairs from Keycloak also so the tokens signed by the compromised key will be invalid.

REST and confidential clients must set **Admin URL** so Keycloak can send clients the pushed not-before policy request.

## [](#_user-storage-federation)Using external storage

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/user-federation.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fuser-federation.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fuser-federation.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Organizations can have databases containing information, passwords, and other credentials. Typically, you cannot migrate existing data storage to a Keycloak deployment so Keycloak can federate existing external user databases. Keycloak supports LDAP and Active Directory, but you can also code extensions for any custom user database by using the Keycloak User Storage SPI.

When a user attempts to log in, Keycloak examines that user’s storage to find that user. If Keycloak does not find the user, Keycloak iterates over each User Storage provider for the realm until it finds a match. Data from the external data storage then maps into a standard user model the Keycloak runtime consumes. This user model then maps to OIDC token claims and SAML assertion attributes.

External user databases rarely have the data necessary to support all the features of Keycloak, so the User Storage Provider can opt to store items locally in Keycloak user data storage. Providers can import users locally and sync periodically with external data storage. This approach depends on the capabilities of the provider and the configuration of the provider. For example, your external user data storage may not support OTP. The OTP can be handled and stored by Keycloak, depending on the provider.

### [](#adding-a-provider)Adding a provider

To add a storage provider, perform the following procedure:

Procedure

1. Click **User Federation** in the menu.
   
   User federation
   
   ![User federation](./images/user-federation.png)
2. Choose to add a **Kerberos** or **LDAP** provider
   
   Keycloak brings you to that provider’s configuration page.

### [](#dealing-with-provider-failures)Dealing with provider failures

If a User Storage Provider fails, you may not be able to log in and view users in the Admin Console. Keycloak does not detect failures when using a Storage Provider to look up a user, so it cancels the invocation. If you have a Storage Provider with a high priority that fails during user lookup, the login or user query fails with an exception and will not fail over to the next configured provider.

Keycloak searches the local Keycloak user database first to resolve users before any LDAP or custom User Storage Provider. Consider creating an administrator account stored in the local Keycloak user database in case of problems connecting to your LDAP and back ends.

Each LDAP and custom User Storage Provider has an `enable` toggle on its Admin Console page. Disabling the User Storage Provider skips the provider when performing queries, so you can view and log in with user accounts in a different provider with lower priority. If your provider uses an `import` strategy and is disabled, imported users are still available for lookup in read-only mode.

When a Storage Provider lookup fails, Keycloak does not fail over because user databases often have duplicate usernames or duplicate emails between them. Duplicate usernames and emails can cause problems because the user loads from one external data store when the admin expects them to load from another data store.

### [](#_ldap)Lightweight Directory Access Protocol (LDAP) and Active Directory

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/user-federation/ldap.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fuser-federation%2Fldap.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fuser-federation%2Fldap.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak includes an LDAP/AD provider. You can federate multiple different LDAP servers in one Keycloak realm and map LDAP user attributes into the Keycloak common user model.

By default, Keycloak maps the username, email, first name, and last name of the user account, but you can also configure additional [mappings](#_ldap_mappers). Keycloak’s LDAP/AD provider supports password validation using LDAP/AD protocols and storage, edit, and synchronization modes.

#### [](#configuring-federated-ldap-storage)Configuring federated LDAP storage

Procedure

1. Click **User Federation** in the menu.
   
   User federation
   
   ![User federation](./images/user-federation.png)
2. Click **Add LDAP providers**.
   
   Keycloak brings you to the LDAP configuration page.
   
   Add LDAP provider
   
   ![User federation](./images/user-fed-ldap.png)

#### [](#storage-mode)Storage mode

Keycloak imports users from LDAP into the local Keycloak user database. This copy of the user database synchronizes on-demand or through a periodic background task. An exception exists for synchronizing passwords. Keycloak never imports passwords. Password validation always occurs on the LDAP server.

The advantage of synchronization is that all Keycloak features work efficiently because any required extra per-user data is stored locally. The disadvantage is that each time Keycloak queries a specific user for the first time, Keycloak performs a corresponding database insert. Also, when imported users are returned as part of a search operation, a corresponding LDAP search is performed for each one to check if the user still exists in LDAP and do some basic validation.

You can synchronize the import with your LDAP server. Import synchronization is unnecessary when LDAP mappers always read particular attributes from the LDAP rather than the database.

You can use LDAP with Keycloak without importing users into the Keycloak user database. The LDAP server backs up the common user model that the Keycloak runtime uses. If LDAP does not support data that a Keycloak feature requires, that feature will not work. The advantage of this approach is that you do not have the resource usage of importing and synchronizing copies of LDAP users into the Keycloak user database.

The **Import Users** switch on the LDAP configuration page controls this storage mode. To import users, toggle this switch to **ON**.

If you disable **Import Users**, you cannot save user profile attributes into the Keycloak database. Also, you cannot save metadata except for user profile metadata mapped to the LDAP. This metadata can include role mappings, group mappings, and other metadata based on the LDAP mappers' configuration.

When you attempt to change the non-LDAP mapped user data, the user update is not possible. For example, you cannot disable the LDAP mapped user unless the user’s `enabled` flag maps to an LDAP attribute.

When working with imported users, Keycloak performs a LDAP search when the user is queried to validate the user and decorate it so that the configured mappers work properly. This means that extra care must be taken when performing unfiltered user searches that may fetch a big number of users as a LDAP search will be issued for every imported user that is found, possibly affecting the performance in a negative way.

Operations that fetch a single user (for example during login) are usually cached and should not be impacted by this extra LDAP search that is performed when the user is fetched for the first time.

#### [](#edit-mode)Edit mode

Users and admins can modify user metadata, users through the [Account Console](#_account-service), and administrators through the Admin Console. The `Edit Mode` configuration on the LDAP configuration page defines the user’s LDAP update privileges.

READ\_ONLY

You cannot change the username, email, first name, last name, and other mapped attributes. Keycloak shows an error anytime a user attempts to update these fields. Password updates are not supported.

WRITABLE

You can change the username, email, first name, last name, and other mapped attributes and passwords and synchronize them automatically with the LDAP store.

UNSYNCED

Keycloak stores changes to the username, email, first name, last name, and passwords in Keycloak local storage, so the administrator must synchronize this data back to LDAP. In this mode, Keycloak deployments can update user metadata on read-only LDAP servers. This option also applies when importing users from LDAP into the local Keycloak user database.

When Keycloak creates the LDAP provider, Keycloak also creates a set of initial [LDAP mappers](#_ldap_mappers). Keycloak configures these mappers based on a combination of the **Vendor**, **Edit Mode**, and **Import Users** switches. For example, when edit mode is UNSYNCED, Keycloak configures the mappers to read a particular user attribute from the database and not from the LDAP server. However, if you later change the edit mode, the mapper’s configuration does not change because it is impossible to detect if the configuration changes changed in UNSYNCED mode. Decide the **Edit Mode** when creating the LDAP provider. This note applies to **Import Users** switch also.

#### [](#other-configuration-options)Other configuration options

Console Display Name

The name of the provider to display in the admin console.

Priority

The priority of the provider when looking up users or adding a user.

Sync Registrations

Toggle this switch to **ON** if you want new users created by Keycloak added to LDAP.

Allow Kerberos authentication

Enable Kerberos/SPNEGO authentication in the realm with user data provisioned from LDAP. For more information, see the [Kerberos section](#_kerberos).

Remove invalid users during searches

Remove users from the local database if they are not available from the user storage when executing searches. If this is true, users no longer available from their corresponding user storage will be deleted from the local database whenever trying to look up users. If false, then users previously imported from the user storage will be kept in the local database, as read-only and disabled, even if that user is no longer available from the user storage. For example, user was deleted directly from LDAP or the `Users DN` is invalid. Note that this behavior will only happen when the user is not yet cached.

Relative User Creation DN

Relative DN from the `Users DN` where new users will be created. This allows users to be created in a sub-DN of the parent `Users DN` when using a `subtree` search scope. For example, if the `Users DN` is set to `ou=people,dc=myorg,dc=com` and the `Relative User Creation DN` is set to `ou=engineering`, users will be fetched from the `Users DN` and all sub-DNs, but new users will be stored in `ou=engineering,ou=people,dc=myorg,dc=com`. In other words, Keycloak concatenates the `Relative User Creation DN` with the `Users DN` (a comma is added automatically when concatenating the DNs) and uses this resulting DN to store users

A similar property is also available in the group and role mappers, allowing groups and roles to be added to a sub-DN of the base DN that is used to search for the groups/roles.

Other options

Hover the mouse pointer over the tooltips in the Admin Console to see more details about these options.

#### [](#connecting-to-ldap-over-ssl)Connecting to LDAP over SSL

When you configure a secure connection URL to your LDAP store (for example,`ldaps://myhost.com:636`), Keycloak uses SSL to communicate with the LDAP server. Configure a truststore on the Keycloak server side so that Keycloak can trust the SSL connection to LDAP - see [Configuring a Truststore](https://www.keycloak.org/server/keycloak-truststore) guide.

The `Use Truststore SPI` configuration property is deprecated. It should normally be left as `Always`.

#### [](#synchronizing-ldap-users-to-keycloak)Synchronizing LDAP users to Keycloak

If you set the **Import Users** option, the LDAP Provider handles importing LDAP users into the Keycloak local database. The first time a user logs in or is returned as part of a user query (e.g. using the search field in the admin console), the LDAP provider imports the LDAP user into the Keycloak database. During authentication, the LDAP password is validated.

By default, Keycloak does not support the username and email attributes with case-sensitive values when storing users to the local database. The value for these attributes will be stored in lower-case in the local database. However, if the **Import Users** option is disabled, Keycloak will not lower-case the username and email attributes when querying users from LDAP. This behavior allows you to use case-sensitive usernames and emails when **Import Users** is disabled. Note that this behavior applies only to username and email attributes. Other attributes remain case-sensitive.

It is recommended to not use case-sensitive usernames and emails when using LDAP with Keycloak, as some features in Keycloak may not work correctly with case-sensitive usernames and emails.

If you want to sync all LDAP users into the Keycloak database, configure and enable the **Sync Settings** on the LDAP provider configuration page.

Two types of synchronization exist:

Periodic Full sync

This type synchronizes all LDAP users into the Keycloak database. The LDAP users already in Keycloak, but different in LDAP, directly update in the Keycloak database.

Periodic Changed users sync

When synchronizing, Keycloak creates or updates users created or updated after the last sync only.

The best way to synchronize is to click **Synchronize all users** when you first create the LDAP provider, then set up periodic synchronization of changed users.

#### [](#_ldap_mappers)LDAP mappers

LDAP mappers are `listeners` triggered by the LDAP Provider. They provide another extension point to LDAP integration. LDAP mappers are triggered when:

- Users log in by using LDAP.
- Users initially register.
- The Admin Console queries a user.

When you create an LDAP Federation provider, Keycloak automatically provides a set of `mappers` for this provider. This set is changeable by users, who can also develop mappers or update/delete existing ones.

User Attribute Mapper

This mapper specifies which LDAP attribute maps to the attribute of the Keycloak user. For example, you can configure the `mail` LDAP attribute to the `email` attribute in the Keycloak database. For this mapper implementation, a one-to-one mapping always exists.

FullName Mapper

This mapper specifies the full name of the user. Keycloak saves the name in an LDAP attribute (usually `cn`) and maps the name to the `firstName` and `lastname` attributes in the Keycloak database. Having `cn` to contain the full name of the user is common for LDAP deployments.

When you register new users in Keycloak and `Sync Registrations` is ON for the LDAP provider, the fullName mapper permits falling back to the username. This fallback is useful when using Microsoft Active Directory (MSAD). The common setup for MSAD is to configure the `cn` LDAP attribute as fullName and, at the same time, use the `cn` LDAP attribute as the `RDN LDAP Attribute` in the LDAP provider configuration. With this setup, Keycloak falls back to the username. For example, if you create Keycloak user "john123" and leave firstName and lastName empty, then the fullname mapper saves "john123" as the value of the `cn` in LDAP. When you enter "John Doe" for firstName and lastName later, the fullname mapper updates LDAP `cn` to the "John Doe" value as falling back to the username is unnecessary.

Hardcoded Attribute Mapper

This mapper adds a hardcoded attribute value to each Keycloak user linked with LDAP. This mapper can also force values for the `enabled` or `emailVerified` user properties.

Role Mapper

This mapper configures role mappings from LDAP into Keycloak role mappings. A single role mapper can map LDAP roles (usually groups from a particular branch of the LDAP tree) into roles corresponding to a specified client’s realm roles or client roles. You can configure more Role mappers for the same LDAP provider. For example, you can specify that role mappings from groups under `ou=main,dc=example,dc=org` map to realm role mappings, and role mappings from groups under `ou=finance,dc=example,dc=org` map to client role mappings of client `finance`.

Hardcoded Role Mapper

This mapper grants a specified Keycloak role to each Keycloak user from the LDAP provider.

Group Mapper

This mapper maps LDAP groups from a branch of an LDAP tree into groups within Keycloak. This mapper also propagates user-group mappings from LDAP into user-group mappings in Keycloak.

MSAD User Account Mapper

This mapper is specific to Microsoft Active Directory (MSAD). It can integrate the MSAD user account state into the Keycloak account state, such as enabled account or expired password. This mapper uses the `userAccountControl`, and `pwdLastSet` LDAP attributes, specific to MSAD and are not the LDAP standard. For example, if the value of `pwdLastSet` is `0`, the Keycloak user must update their password. The result is an UPDATE\_PASSWORD required action added to the user. If the value of `userAccountControl` is `514` (disabled account), the Keycloak user is disabled.

Certificate Mapper

This mapper maps X.509 certificates. Keycloak uses it in conjunction with X.509 authentication and `Full certificate in PEM format` as an identity source. This mapper behaves similarly to the `User Attribute Mapper`, but Keycloak can filter for an LDAP attribute storing a PEM or DER format certificate. Enable `Always Read Value From LDAP` with this mapper.

User Attribute mappers that map basic Keycloak user attributes, such as username, firstname, lastname, and email, to corresponding LDAP attributes. You can extend these and provide your own additional attribute mappings. The Admin Console provides tooltips to help with configuring the corresponding mappers.

#### [](#_ldap_password_hashing)Password hashing

When Keycloak updates a password, Keycloak sends the password in plain-text format. This action is different from updating the password in the built-in Keycloak database, where Keycloak hashes and salts the password before sending it to the database. For LDAP, Keycloak relies on the LDAP server to hash and salt the password.

By default, LDAP servers such as MSAD, RHDS, or FreeIPA hash and salt passwords. Other LDAP servers such as OpenLDAP store the passwords in plain-text unless you use the *LDAPv3 Password Modify Extended Operation* as described in [RFC 3062](https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.3). Enable the LDAPv3 Password Modify Extended Operation in the LDAP configuration page. See the documentation of your LDAP server for more details. [Configure ApacheDS to hash and salt passwords automatically](https://directory.apache.org/apacheds/advanced-ug/4.1.1.4-ss-password-hash.html) by enabling the passwordHashing interceptor.

Always verify that user passwords are properly hashed and not stored as plaintext by inspecting a changed directory entry using `ldapsearch` and base64 decode the `userPassword` attribute value.

#### [](#_ldap_password_policy)Enabling password change after reset

You can force users to change their password after an administrator resets their password. This feature is useful for security reasons because it prevents users from using the temporary password set by the administrator for a long time.

For that, enable the `Enable LDAP password policy` setting so that a `UPDATE_PASSWORD` required action will be added to the user whenever the LDAP server indicates that the user must update their password.

Keycloak is usually configured to connect to the LDAP server using an administrator account even when users are updating their password. In this case, some LDAP servers, if not configured properly, will always force a password change after the user changes their password because the LDAP server will see that the password was changed by an administrator account and not by the user itself.

When enabling this feature, make sure your LDAP server is configured to not force users to change their password a second time after the user changes their password. Not doing so might cause a password change loop where the user is forced to change their password every time they log in.

This capability is based on [Password Policy for LDAP Directories (IETF draft-behera-ldap-password-policy)](https://datatracker.ietf.org/doc/html/draft-behera-ldap-password-policy-11). For Microsoft Active Directory, use the [MSAD User Account Mapper](#_msad_mapper) instead.

#### [](#_ldap_connection_pool)Configuring the connection pool

For more efficiency when managing LDAP connections and to improve performance when handling multiple connections, you can enable connection pooling. By doing that, when a connection is closed, it will be returned to the pool for future use therefore reducing the cost of creating new connections all the time.

The LDAP connection pool configuration is configured using the following system properties:

  Name Description

`com.sun.jndi.ldap.connect.pool.authentication`

A list of space-separated authentication types of connections that may be pooled. Valid types are "none", "simple", and "DIGEST-MD5"

`com.sun.jndi.ldap.connect.pool.initsize`

The string representation of an integer that represents the number of connections per connection identity to create when initially creating a connection for the identity

`com.sun.jndi.ldap.connect.pool.maxsize`

The string representation of an integer that represents the maximum number of connections per connection identity that can be maintained concurrently. Note setting this value too low may cause contention in obtaining LDAP connections. See also `com.sun.jndi.ldap.connect.timeout`.

`com.sun.jndi.ldap.connect.pool.prefsize`

The string representation of an integer that represents the preferred number of connections per connection identity that should be maintained concurrently

`com.sun.jndi.ldap.connect.pool.timeout`

The string representation of an integer that represents the number of milliseconds that an idle connection may remain in the pool without being closed and removed from the pool

`com.sun.jndi.ldap.connect.pool.protocol`

A list of space-separated protocol types of connections that may be pooled. Valid types are "plain" and "ssl"

`com.sun.jndi.ldap.connect.pool.debug`

A string that indicates the level of debug output to produce. Valid values are "fine" (trace connection creation and removal) and "all" (all debugging information)

`com.sun.jndi.ldap.connect.timeout`

The string representation of an integer that represents how long in milliseconds obtaining a connection should take. This is also applicable to wait times due to connection pool contention. Effectively defaults to 5000.

By default, connection pooling is enabled for both `plain` and `ssl` protocols.

For more details, see the [Java LDAP Connection Pooling Configuration](https://docs.oracle.com/javase/jndi/tutorial/ldap/connect/config.html) documentation.

To set any of these properties, you can set the `JAVA_OPTS_APPEND` environment variable:

```
export JAVA_OPTS_APPEND=-Dcom.sun.jndi.ldap.connect.pool.initsize=10 -Dcom.sun.jndi.ldap.connect.pool.maxsize=50
```

#### [](#_ldap_troubleshooting)Troubleshooting

It is useful to increase the logging level to TRACE for the category `org.keycloak.storage.ldap`. With this setting, many logging messages are sent to the server log in the `TRACE` level, including the logging for all queries to the LDAP server and the parameters, which were used to send the queries. When you are creating any LDAP question on user forum or JIRA, consider attaching the server log with enabled TRACE logging. If it is too big, the good alternative is to include just the snippet from server log with the messages, which were added to the log during the operation, which causes the issues to you.

- When you create an LDAP provider, a message appears in the server log in the INFO level starting with:

```
Creating new LDAP Store for the LDAP storage provider: ...
```

It shows the configuration of your LDAP provider. Before you are asking the questions or reporting bugs, it will be nice to include this message to show your LDAP configuration. Eventually feel free to replace some config changes, which you do not want to include, with some placeholder values. One example is `bindDn=some-placeholder` . For `connectionUrl`, feel free to replace it as well, but it is generally useful to include at least the protocol, which was used (`ldap` vs `ldaps`)\`. Similarly it can be useful to include the details for configuration of your LDAP mappers, which are displayed with the message like this at the DEBUG level:

```
Mapper for provider: XXX, Mapper name: YYY, Provider: ZZZ ...
```

Note those messages are displayed just with the enabled DEBUG logging.

- For tracking the performance or connection pooling issues, consider setting the value of property `com.sun.jndi.ldap.connect.pool.debug` to `all`. This change adds many additional messages to the server log with the included logging for the LDAP connection pooling. As a result, you can track the issues related to connection pooling or performance. For more details, see [Configuring the connection pool](#_ldap_connection_pool).

After changing the configuration of connection pooling, you may need to restart the Keycloak server to enforce re-initialization of the LDAP provider connection.

If no more messages appear for connection pooling even after server restart, it can indicate that connection pooling does not work with your LDAP server.

- For the case of reporting LDAP issue, you may consider to attach some part of your LDAP tree with the target data, which causes issues in your environment. For example if login of some user takes lot of time, you can consider attach his LDAP entry showing count of `member` attributes of various "group" entries. In this case, it might be useful to add if those group entries are mapped to some Group LDAP mapper (or Role LDAP Mapper) in Keycloak and so on.

### [](#_sssd)SSSD and FreeIPA Identity Management integration

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/user-federation/sssd.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fuser-federation%2Fsssd.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fuser-federation%2Fsssd.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak includes the [System Security Services Daemon (SSSD)](https://fedoraproject.org/wiki/Features/SSSD) plugin. SSSD is part of the Fedora and Red Hat Enterprise Linux (RHEL), and it provides access to multiple identities and authentication providers. SSSD also provides benefits such as failover and offline support. For more information, see [the Red Hat Enterprise Linux Identity Management documentation](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/7/html/system-level_authentication_guide/sssd).

SSSD integrates with the FreeIPA identity management (IdM) server, providing authentication and access control. With this integration, Keycloak can authenticate against privileged access management (PAM) services and retrieve user data from SSSD. For more information about using Red Hat Identity Management in Linux environments, see [the Red Hat Enterprise Linux Identity Management documentation](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/7/html/linux_domain_identity_authentication_and_policy_guide/index).

![keycloak sssd freeipa integration overview](./images/keycloak-sssd-freeipa-integration-overview.png)

Keycloak and SSSD communicate through read-only D-Bus interfaces. For this reason, the way to provision and update users is to use the FreeIPA/IdM administration interface. By default, the interface imports the username, email, first name, and last name.

Keycloak registers groups and roles automatically but does not synchronize them. The groups are imported from SSSD the first time the user is accessed and then they are managed entirely inside Keycloak. Any changes made by the administrator in Keycloak do not synchronize with SSSD or vice-versa.

#### [](#freeipaidm-server)FreeIPA/IdM server

The [FreeIPA Container image](https://quay.io/repository/freeipa/freeipa-server?tab=tags%2F) is available at [Quay.io](https://quay.io/). To set up the FreeIPA server, see the [FreeIPA documentation](https://www.freeipa.org/page/Quick_Start_Guide).

Procedure

1. Run your FreeIPA server using this command:
   
   ```
    docker run --name freeipa-server-container -it \
    -h server.freeipa.local -e PASSWORD=YOUR_PASSWORD \
    -v /sys/fs/cgroup:/sys/fs/cgroup:ro \
    -v /var/lib/ipa-data:/data:Z freeipa/freeipa-server
   ```
   
   The parameter `-h` with `server.freeipa.local` represents the FreeIPA/IdM server hostname. Change `YOUR_PASSWORD` to a password of your own.
2. After the container starts, change the `/etc/hosts` file to include:
   
   ```
   x.x.x.x     server.freeipa.local
   ```
   
   If you do not make this change, you must set up a DNS server.
3. Use the following command to enroll your Linux server in the IPA domain so that the SSSD federation provider starts and runs on Keycloak:
   
   ```
    ipa-client-install --mkhomedir -p admin -w password
   ```
4. Run the following command on the client to verify the installation is working:
   
   ```
    kinit admin
   ```
5. Enter your password.
6. Add users to the IPA server using this command:
   
   ```
   $ ipa user-add <username> --first=<first name> --last=<surname> --email=<email address> --phone=<telephoneNumber> --street=<street> --city=<city> --state=<state> --postalcode=<postal code> --password
   ```
7. Force set the user’s password using kinit.
   
   ```
    kinit <username>
   ```
8. Enter the following to restore normal IPA operation:
   
   ```
   kdestroy -A
   kinit admin
   ```

#### [](#sssd-and-d-bus)SSSD and D-Bus

The federation provider obtains the data from SSSD using D-BUS. It authenticates the data using PAM.

Procedure

1. Install the sssd-dbus RPM.
   
   ```
   $ sudo yum install sssd-dbus
   ```
2. Run the following provisioning script:
   
   ```
   $ bin/federation-sssd-setup.sh
   ```
   
   The script can also be used as a guide to configure SSSD and PAM for Keycloak. It makes the following changes to `/etc/sssd/sssd.conf`:
   
   ```
     [domain/your-hostname.local]
     ...
     ldap_user_extra_attrs = mail:mail, sn:sn, givenname:givenname, telephoneNumber:telephoneNumber
     ...
     [sssd]
     services = nss, sudo, pam, ssh, ifp
     ...
     [ifp]
     allowed_uids = root, yourOSUsername
     user_attributes = +mail, +telephoneNumber, +givenname, +sn
   ```
   
   The `ifp` service is added to SSSD and configured to allow the OS user to interrogate the IPA server through this interface.
   
   The script also creates a new PAM service `/etc/pam.d/keycloak` to authenticate users via SSSD:
   
   ```
   auth    required   pam_sss.so
   account required   pam_sss.so
   ```
3. Run `dbus-send` to ensure the setup is successful.
   
   ```
   dbus-send --print-reply --system --dest=org.freedesktop.sssd.infopipe /org/freedesktop/sssd/infopipe org.freedesktop.sssd.infopipe.GetUserAttr string:<username> array:string:mail,givenname,sn,telephoneNumber
   
   dbus-send --print-reply --system --dest=org.freedesktop.sssd.infopipe /org/freedesktop/sssd/infopipe org.freedesktop.sssd.infopipe.GetUserGroups string:<username>
   ```
   
   If the setup is successful, each command displays the user’s attributes and groups respectively. If there is a timeout or an error, the federation provider running on Keycloak cannot retrieve any data. This error usually happens because the server is not enrolled in the FreeIPA IdM server, or does not have permission to access the SSSD service.
   
   If you do not have permission to access the SSSD service, ensure that the user running the Keycloak server is in the `/etc/sssd/sssd.conf` file in the following section:
   
   ```
   [ifp]
   allowed_uids = root, yourOSUsername
   ```
   
   And the `ipaapi` system user is created inside the host. This user is necessary for the `ifp` service. Check the user is created in the system.
   
   ```
   grep ipaapi /etc/passwd
   ipaapi:x:992:988:IPA Framework User:/:/sbin/nologin
   ```

#### [](#enabling-the-sssd-federation-provider)Enabling the SSSD federation provider

Keycloak uses [DBus-Java](https://github.com/hypfvieh/dbus-java) project to communicate at a low level with D-Bus and [JNA](https://github.com/java-native-access/jna) to authenticate via Operating System Pluggable Authentication Modules (PAM).

Although now Keycloak contains all the needed libraries to run the `SSSD` provider, JDK version 21 or later is needed. Therefore the `SSSD` provider will only be displayed when the host configuration is correct and JDK 21 or later is used to run Keycloak.

#### [](#configuring-a-federated-sssd-store)Configuring a federated SSSD store

After the installation, configure a federated SSSD store.

Procedure

1. Click **User Federation** in the menu.
2. If everything is setup successfully the **Add Sssd providers** button will be displayed in the page. Click on it.
3. Assign a name to the new provider.
4. Click **Save**.

You can now authenticate against Keycloak using a FreeIPA/IdM user and credentials.

### [](#custom-providers)Custom providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/user-federation/custom.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fuser-federation%2Fcustom.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fuser-federation%2Fcustom.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak does have a Service Provider Interface (SPI) for User Storage Federation to develop custom providers. You can find documentation on developing customer providers in the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/).

## [](#assembly-managing-users_server_administration_guide)Managing users

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/assembly-managing-users.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fassembly-managing-users.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fassembly-managing-users.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

From the Admin Console, you have a wide range of actions you can perform to manage users.

### [](#proc-creating-user_server_administration_guide)Creating users

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-creating-user.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-creating-user.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-creating-user.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You create users in the realm where you intend to have applications needed by those users. Avoid creating users in the master realm, which is only intended for creating other realms.

Prerequisite

- You are in a realm other than the master realm.

Procedure

1. Click **Users** in the menu.
2. Click **Add User**.
3. Enter the details for the new user.
   
   **Username** is the only required field.
4. Click **Save**. After saving the details, the **Management** page for the new user is displayed.

### [](#user-profile)Managing user attributes

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/user-profile.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fuser-profile.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fuser-profile.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

In Keycloak a user is associated with a set of attributes. These attributes are used to better describe and identify users within Keycloak as well as to pass over additional information about them to applications.

A user profile defines a well-defined schema for representing user attributes and how they are managed within a realm. By providing a consistent view over user information, it allows administrators to control the different aspects on how attributes are managed as well as to make it much easier to extend Keycloak to support additional attributes.

Although the user profile is mainly targeted for attributes that end-users can manage (e.g.: first and last names, phone, etc) it also serves for managing any other metadata you want to associate with your users.

Among other capabilities, user profile enables administrators to:

- Define a schema for user attributes
- Define whether an attribute is required based on contextual information (e.g.: if required only for users, or admins, or both, or depending on the scope being requested.)
- Define specific permissions for viewing and editing user attributes, making possible to adhere to strong privacy requirements where some attributes can not be seen or be changed by third-parties (including administrators)
- Dynamically enforce user profile compliance so that user information is always updated and in compliance with the metadata and rules associated with attributes
- Define validation rules on a per-attribute basis by leveraging the built-in validators or writing custom ones
- Dynamically render forms that users interact with like registration, update profile, brokering, and personal information in the account console, according to the attribute definitions and without any need to manually change themes.
- Customize user management interfaces in the administration console so that attributes are rendered dynamically based on the user profile schema

The user profile schema or configuration uses a [JSON](#_user-profile-json-configuration) format to represent attributes and their metadata. From the administration console, you are able to manage the configuration by clicking on the `Realm Settings` on the left side menu and then clicking on the `User Profile` tab on that page.

In the next sections, we’ll be looking at how to create your own user profile schema or configuration, and how to manage attributes.

#### [](#understanding-the-default-configuration)Understanding the Default Configuration

By default, Keycloak provides a basic user profile configuration covering some of the most common user attributes:

  Name Description

`username`

The username

`email`

End-User’s preferred e-mail address.

`firstName`

Given name(s) or first name(s) of the end-user

`lastName`

Surname(s) or last name(s) of the End-User

In Keycloak, both `username` and `email` attributes have a special handling as they are often used to identify, authenticate, and link user accounts. For those attributes, you are limited to changing their settings, and you can not remove them.

The behavior of both `username` and `email` attributes changes accordingly to the `Login` settings of your realm. For instance, changing the `Email as username` or the `Edit username` settings will override any configuration you have set in the user profile configuration.

As you will see in the following sections, you are free to change the default configuration by bringing your own attributes or changing the settings for any of the available attributes to better fit it to your needs.

#### [](#understanding-the-user-profile-contexts)Understanding the User Profile Contexts

In Keycloak, users are managed through different contexts:

- Registration
- Update Profile
- Reviewing Profile when authenticating through a broker or social provider
- Account Console
- Administrative (e.g.: administration console and Admin REST API)

Except for the `Administrative` context, all other contexts are considered end-user contexts as they are related to user self-service flows.

Knowing these contexts is important to understand where your user profile configuration will take effect when managing users. Regardless of the context where the user is being managed, the same user profile configuration will be used to render UIs and validate attribute values.

As you will see in the following sections, you might restrict certain attributes to be available only from the administrative context and disable them completely for end-users. The other way around is also true if you don’t want administrators to have access to certain user attributes but only the end-user.

#### [](#_understanding-managed-and-unmanaged-attributes)Understanding Managed and Unmanaged Attributes

By default, Keycloak will only recognize the attributes defined in your user profile configuration. The server ignores any other attribute not explicitly defined there.

By being strict about which user attributes can be set to your users, as well as how their values are validated, Keycloak can add another defense barrier to your realm and help you to prevent unexpected attributes and values associated to your users.

That said, user attributes can be categorized as follows:

- **Managed**. These are attributes controlled by your user profile, to which you want to allow end-users and administrators to manage from any user profile context. For these attributes, you want complete control on how and when they are managed.
- **Unmanaged**. These are attributes you do not explicitly define in your user profile so that they are completely ignored by Keycloak, by default.

Although unmanaged attributes are disabled by default, you can configure your realm using different policies to define how they are handled by the server. For that, click on the `Realm Settings` at the left side menu, click on the `General` tab, and then choose any of the following options from the `Unmanaged Attributes` setting:

- **Disabled**. This is the default policy so that unmanaged attributes are disabled from all user profile contexts.
- **Enabled**. This policy enables unmanaged attributes to all user profile contexts.
- **Admin can view**. This policy enables unmanaged attributes only from the administrative context as read-only.
- **Admin can edit**. This policy enables unmanaged attributes only from the administrative context for reads and writes.

These policies give you a fine-grained control over how the server will handle unmanaged attributes. You can choose to completely disable or only support unmanaged attributes when managing users through the administrative context.

When unmanaged attributes are enabled (even if partially) you can manage them from the administration console at the `Attributes` tab in the User Details UI. If the policy is set to `Disabled` this tab is not available.

As a security recommendation, try to adhere to the most strict policy as much as possible (e.g.: `Disabled` or `Admin can edit`) to prevent unexpected attributes (and values) set to your users when they are managing their profile through end-user contexts. Avoid setting the `Enabled` policy and prefer defining all the attributes that end-users can manage in your user profile configuration, under your control.

The `Enabled` policy is targeted for realms migrating from previous versions of Keycloak and to avoid breaking behavior when using custom themes and extending the server with their own custom user attributes.

As you will see in the following sections, you can also restrict the audience for an attribute by choosing if it should be visible or writable by users and/or administrators.

For unmanaged attributes, the maximum length is 2048 characters. To specify a different minimum or maximum length, change the unmanaged attribute to a managed attribute and add a `length` validator.

Keycloak caches user-related objects in its internal caches. The longer the attributes are, the more memory the cache consumes. Therefore, limiting the size of the length attributes is recommended. Consider storing large objects outside Keycloak and reference them by ID or URL.

#### [](#managing-the-user-profile)Managing the User Profile

The user profile configuration is managed on a per-realm basis. For that, click on the `Realm Settings` link on the left side menu and then click on the `User Profile` tab.

User Profile Tab

![user profile tab](./images/user-profile-tab.png)

In the `Attributes` sub-tab you have a list of all managed attributes.

In the `Attribute Groups` sub-tab you can manage attribute groups. An attribute group allows you to correlate attributes so that they are displayed together when rendering user facing forms.

In the `JSON Editor` sub-tab you can view and edit the [JSON](#_user-profile-json-configuration) configuration. You can use this tab to grab your current configuration or manage it manually. Any change you make to this tab is reflected in the other tabs, and vice-versa.

In the next section, you are going to learn how to manage attributes.

#### [](#managing-attributes)Managing Attributes

At the `Attributes` sub-tab you can create, edit, and delete the managed attributes.

To define a new attribute and associate it with the user profile, click on the **Create attribute** button at the top of the attribute listing.

Attribute Configuration

![user profile create attribute](./images/user-profile-create-attribute.png)

When configuring the attribute you can define the following settings:

Name

The name of the attribute, used to uniquely identify an attribute.

Display name

A user-friendly name for the attribute, mainly used when rendering user-facing forms. It also supports [Using Internationalized Messages](#_using-internationalized-messages).

Multivalued

If enabled, the attribute supports multiple values and UIs are rendered accordingly to allow setting many values. When enabling this setting, make sure to add a validator to set a hard limit to the number of values.

Default Value

Defines the value that will be automatically assigned to the attribute if the user does not provide one. Ensure that this default value complies with all configured validators for the attribute.

Attribute Group

The attribute group to which the attribute belongs to, if any.

Enabled when

Enables or disables an attribute. If set to `Always`, the attribute is available from any user profile context. If set to `Scopes are requested`, the attribute is only available when the client acting on behalf of the user is requesting a set of one or more scopes. You can use this option to dynamically enforce certain attributes depending on the client scopes being requested. For the administration console, scopes are not evaluated and the attribute is always enabled. That is because filtering attributes by scopes only works when running end-user authentication flows.

Required

Set the conditions to mark an attribute as required. If disabled, the attribute is optional. If enabled, you can set the `Required for` setting to mark the attribute as required depending on the user profile context so that the attribute is required for end-users (via end-user contexts) or to administrators (via administrative context), or both. You can also set the `Required when` setting to mark the attribute as required only when a set of one or more client scopes are requested. If set to `Always`, the attribute is required from any user profile context. If set to `Scopes are requested`, the attribute is only required when the client acting on behalf of the user is requesting a set of one or more scopes. For the account and administration consoles, scopes are not evaluated and the attribute is not required. That is because filtering attributes by scopes only works when running authentication flows.

Permission

In this section, you can define read and write permissions when the attribute is being managed from an end-user or administrative context. The `Who can edit` setting mark an attribute as writable by `User` and/or `Admin`, from an end-user and administrative context, respectively. The `Who can view` setting mark an attribute as read-only by `User` and/or `Admin` from an end-user and administrative context, respectively.

Validation

In this section, you can define the validations that will be performed when managing the attribute value. Keycloak provides a set of built-in validators you can choose from with the possibility to add your own. For more details, look at the [Validating Attributes](#_validating-attributes) section.

Annotation

In this section, you can associate annotations to the attribute. Annotations are mainly useful to pass over additional metadata to frontends for rendering purposes. For more details, look at the [Defining UI Annotations](#_defining-ui-annotations) section.

When you create an attribute, the attribute is only available from administrative contexts to avoid unexpectedly exposing attributes to end-users. Effectively, the attribute won’t be accessible to end-users when they are managing their profile through the end-user contexts. You can change the `Permissions` settings anytime accordingly to your needs.

#### [](#_validating-attributes)Validating Attributes

You can enable validation to managed attributes to make sure the attribute value conforms to specific rules. For that, you can add or remove validators from the `Validations` settings when managing an attribute.

Attribute Validation

![user profile validation](./images/user-profile-validation.png)

Validation happens at any time when writing to an attribute, and they can throw errors that will be shown in UIs when the value fails a validation.

For security reasons, every attribute that is editable by users should have a validation to restrict the size of the values users enter. If no `length` validator has been specified, Keycloak defaults to a maximum length of 2048 characters.

##### [](#built-in-validators)Built-in Validators

Keycloak provides some built-in validators that you can choose from, and you are also able to provide your own validators by extending the `Validator SPI`.

The list below provides a list of all the built-in validators:

   Name Description Configuration

length

Check the length of a string value based on a minimum and maximum length.

**min**: an integer to define the minimum allowed length.

**max**: an integer to define the maximum allowed length.

**trim-disabled**: a boolean to define whether the value is trimmed prior to validation.

integer

Check if the value is an integer and within a lower and/or upper range. If no range is defined, the validator only checks whether the value is a valid number.

**min**: an integer to define the lower range.

**max**: an integer to define the upper range.

double

Check if the value is a double and within a lower and/or upper range. If no range is defined, the validator only checks whether the value is a valid number.

**min**: an integer to define the lower range.

**max**: an integer to define the upper range.

uri

Check if the value is a valid URI.

None

pattern

Check if the value matches a specific RegEx pattern.

**pattern**: the RegEx pattern to use when validating values.

**error-message**: the key of the error message in i18n bundle. If not set a generic message is used.

email

Check if the value has a valid e-mail format. If [the realm is configured to send emails](#_email) and the option **Allow UTF-8** is not enabled to support internationalized emails, this validator also checks that the local part of the address contains only ASCII characters.

**max-local-length**: an integer to define the maximum length for the local part of the email. It defaults to 64 per specification.

local-date

Check if the value has a valid format based on the realm and/or user locale.

None

iso-date

Check if the value has a valid format based on ISO 8601. This validator can be used with inputs using the html5-date input type.

None

person-name-prohibited-characters

Check if the value is a valid person name as an additional barrier for attacks such as script injection. The validation is based on a default RegEx pattern that blocks characters not common in person names.

**error-message**: the key of the error message in i18n bundle. If not set a generic message is used.

username-prohibited-characters

Check if the value is a valid username as an additional barrier for attacks such as script injection. The validation is based on a default RegEx pattern that blocks characters not common in usernames. When the realm setting `Email as username` is enabled, this validator is skipped to allow email values.

**error-message**: the key of the error message in i18n bundle. If not set a generic message is used.

options

Check if the value is from the defined set of allowed values. Useful to validate values entered through select and multiselect fields.

**options**: array of strings containing allowed values.

up-username-not-idn-homograph

The field can contain only latin characters and common unicode characters. Useful for the fields, which can be subject of IDN homograph attacks (typically username).

**error-message**: the key of the error message in i18n bundle. If not set a generic message is used.

multivalued

Validates the size of a multivalued attribute.

**min**: an integer to define the minimum allowed count of attribute values.

**max**: an integer to define the maximum allowed count of attribute values.

#### [](#_defining-ui-annotations)Defining UI Annotations

In order to pass additional information to frontends, attributes can be decorated with annotations to dictate how attributes are rendered. This capability is mainly useful when extending Keycloak themes to render pages dynamically based on the annotations associated with attributes.

Annotations are used, for example, for [Changing the HTML `type` for an Attribute](#_changing-the-html-type-for-an-attribute) and [Changing the DOM representation of an Attribute](#_changing-the-dom-representation-of-an-attribute), as you will see in the following sections.

Attribute Annotation

![user profile annotation](./images/user-profile-annotation.png)

An annotation is a key/value pair shared with the UI so that they can change how the HTML element corresponding to the attribute is rendered. You can set any annotation you want to an attribute as long as the annotation is supported by the theme your realm is using.

The only restriction you have is to avoid using annotations using the `kc` prefix in their keys because these annotations using this prefix are reserved for Keycloak.

##### [](#built-in-annotations)Built-in Annotations

The following annotations are supported by Keycloak built-in themes:

  Name Description

inputType

Type of the form input field. Available types are described in a table below.

inputHelperTextBefore

Helper text rendered before (above) the input field. Direct text or internationalization pattern (like `${i18n.key}`) can be used here. Text is NOT html escaped when rendered into the page, so you can use html tags here to format the text, but you also have to correctly escape html control characters.

inputHelperTextAfter

Helper text rendered after (under) the input field. Direct text or internationalization pattern (like `${i18n.key}`) can be used here. Text is NOT html escaped when rendered into the page, so you can use html tags here to format the text, but you also have to correctly escape html control characters.

inputOptionsFromValidation

Annotation for select and multiselect types. Optional name of custom attribute validation to get input options from. See the [detailed description](#_managing_options_for_select_fields) below.

inputOptionLabelsI18nPrefix

Annotation for select and multiselect types. Internationalization key prefix to render options in UI. See the [detailed description](#_managing_options_for_select_fields) below.

inputOptionLabels

Annotation for select and multiselect types. Optional map to define UI labels for options (directly or using internationalization). See the [detailed description](#_managing_options_for_select_fields) below.

inputTypePlaceholder

HTML input `placeholder` attribute applied to the field - specifies a short hint that describes the expected value of an input field (e.g. a sample value or a short description of the expected format). The short hint is displayed in the input field before the user enters a value.

inputTypeSize

HTML input `size` attribute applied to the field - specifies the width, in characters, of a single line input field. For fields based on HTML `select` type it specifies number of rows with options shown. May not work, depending on css in used theme!

inputTypeCols

HTML input `cols` attribute applied to the field - specifies the width, in characters, for `textarea` type. May not work, depending on css in used theme!

inputTypeRows

HTML input `rows` attribute applied to the field - specifies the height, in characters, for `textarea` type. For select fields it specifies number of rows with options shown. May not work, depending on css in used theme!

inputTypePattern

HTML input `pattern` attribute applied to the field providing client side validation - specifies a regular expression that an input field’s value is checked against. Useful for single line inputs.

inputTypeMaxLength

HTML input `maxlength` attribute applied to the field providing client side validation - maximal length of the text which can be entered into the input field. Useful for text fields.

inputTypeMinLength

HTML input `minlength` attribute applied to the field providing client side validation - minimal length of the text which can be entered into the input field. Useful for text fields.

inputTypeMax

HTML input `max` attribute applied to the field providing client side validation - maximal value which can be entered into the input field. Useful for numeric fields.

inputTypeMin

HTML input `min` attribute applied to the field providing client side validation - minimal value which can be entered into the input field. Useful for numeric fields.

inputTypeStep

HTML input `step` attribute applied to the field - Specifies the interval between legal numbers in an input field. Useful for numeric fields.

Number Format

If set, the `data-kcNumberFormat` attribute is added to the field to format the value based on a given format. This annotation is targeted for numbers where the format is based on the number of digits expected in a determined position. For instance, a format `({2}) {5}-{4}` will format the field value to `(00) 00000-0000`.

Number UnFormat

If set, the `data-kcNumberUnFormat` attribute is added to the field to format the value based on a given format before submitting the form. This annotation is useful if you do not want to store any format for a specific attribute but only format the value on the client side. For instance, if the current value is `(00) 00000-0000`, the value will change to `00000000000` if you set the value `{11}` to this annotation or any other format you want by specifying a set of one or more group of digits. Make sure to add validators to perform server-side validations before storing values.

Field types use HTML form field tags and attributes applied to them - they behave based on the HTML specifications and browser support for them.

Visual rendering also depends on css styles applied in the used theme.

##### [](#_changing-the-html-type-for-an-attribute)Changing the HTML `type` for an Attribute

You can change the `type` of a HTML5 input element by setting the `inputType` annotation. The available types are:

   Name Description HTML tag used

text

Single line text input.

input

textarea

Multiple line text input.

textarea

select

Common single select input. See the [description of how to configure options](#_managing_options_for_select_fields) below.

select

select-radiobuttons

Single select input through group of radio buttons. See the [description of how to configure options](#_managing_options_for_select_fields) below.

group of input

multiselect

Common multiselect input. See the [description of how to configure options](#_managing_options_for_select_fields) below.

select

multiselect-checkboxes

Multiselect input through group of checkboxes. See the [description of how to configure options](#_managing_options_for_select_fields) below.

group of input

html5-email

Single line text input for email address based on HTML 5 spec.

input

html5-tel

Single line text input for phone number based on HTML 5 spec.

input

html5-url

Single line text input for URL based on HTML 5 spec.

input

html5-number

Single line input for number (integer or float depending on `step`) based on HTML 5 spec.

input

html5-range

Slider for number entering based on HTML 5 spec.

input

html5-datetime-local

Date Time input based on HTML 5 spec.

input

html5-date

Date input based on HTML 5 spec.

input

html5-month

Month input based on HTML 5 spec.

input

html5-week

Week input based on HTML 5 spec.

input

html5-time

Time input based on HTML 5 spec.

input

##### [](#_managing_options_for_select_fields)Defining options for select and multiselect fields

Options for select and multiselect fields are taken from validation applied to the attribute to be sure validation and field options presented in UI are always consistent. By default, options are taken from built-in `options` validation.

You can use various ways to provide nice human-readable labels for select and multiselect options. The simplest case is when attribute values are same as UI labels. No extra configuration is necessary in this case.

Option values same as UI labels

![user profile select options simple](./images/user-profile-select-options-simple.png)

When attribute value is kind of ID not suitable for UI, you can use simple internationalization support provided by `inputOptionLabelsI18nPrefix` annotation. It defines prefix for internationalization keys, option value is dot appended to this prefix.

Simple internationalization for UI labels using i18n key prefix

![user profile select options simple i18n](./images/user-profile-select-options-simple-i18n.png)

Localized UI label texts for option value have to be provided by `userprofile.jobtitle.sweng` and `userprofile.jobtitle.swarch` keys then, using common localization mechanism.

You can also use `inputOptionLabels` annotation to provide labels for individual options. It contains a map of labels for option - key in the map is option value (defined in validation), and value in the map is UI label text itself or its internationalization pattern (like `${i18n.key}`) for that option.

You have to use User Profile `JSON Editor` to enter map as `inputOptionLabels` annotation value.

Example of directly entered labels for individual options without internationalization:

```
"attributes": [
<...
{
  "name": "jobTitle",
  "validations": {
    "options": {
      "options":[
        "sweng",
        "swarch"
      ]
    }
  },
  "annotations": {
    "inputType": "select",
    "inputOptionLabels": {
      "sweng": "Software Engineer",
      "swarch": "Software Architect"
    }
  }
}
...
]
```

Example of the internationalized labels for individual options:

```
"attributes": [
...
{
  "name": "jobTitle",
  "validations": {
    "options": {
      "options":[
        "sweng",
        "swarch"
      ]
    }
  },
  "annotations": {
    "inputType": "select-radiobuttons",
    "inputOptionLabels": {
      "sweng": "${jobtitle.swengineer}",
      "swarch": "${jobtitle.swarchitect}"
    }
  }
}
...
]
```

Localized texts have to be provided by `jobtitle.swengineer` and `jobtitle.swarchitect` keys then, using common localization mechanism.

Custom validator can be used to provide options thanks to `inputOptionsFromValidation` attribute annotation. This validation have to have `options` config providing array of options. Internationalization works the same way as for options provided by built-in `options` validation.

Options provided by custom validator

![user profile select options custom validator](./images/user-profile-select-options-custom-validator.png)

##### [](#_changing-the-dom-representation-of-an-attribute)Changing the DOM representation of an Attribute

You can enable additional client-side behavior by setting annotations with the `kc` prefix. These annotations are going to translate into an HTML attribute in the corresponding element of an attribute, prefixed with `data-`, and a script with the same name will be loaded to the dynamic pages so that you can select elements from the DOM based on the custom `data-` attribute and decorate them accordingly by modifying their DOM representation.

For instance, if you add a `kcMyCustomValidation` annotation to an attribute, the HTML attribute `data-kcMyCustomValidation` is added to the corresponding HTML element for the attribute, and a JavaScript module is loaded from your custom theme at `<THEME TYPE>/resources/js/kcMyCustomValidation.js`. See the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more information about how to deploy a custom JavaScript module to your theme.

The JavaScript module can run any code to customize the DOM and the elements rendered for each attribute. For that, you can use the `userProfile.js` module to register an annotation descriptor for your custom annotation as follows:

```
import { registerElementAnnotatedBy } from "./userProfile.js";

registerElementAnnotatedBy({
  name: 'kcMyCustomValidation',
  onAdd(element) {
    var listener = function (event) {
        // do something on keyup
    };

    element.addEventListener("keyup", listener);

    // returns a cleanup function to remove the event listener
    return () => element.removeEventListener("keyup", listener);
  }
});
```

The `registerElementAnnotatedBy` is a method to register annotation descriptors. A descriptor is an object with a `name`, referencing the annotation name, and a `onAdd` function. Whenever the page is rendered or an attribute with the annotation is added to the DOM, the `onAdd` function is invoked so that you can customize the behavior for the element.

The `onAdd` function can also return a function to perform a cleanup. For instance, if you are adding event listeners to elements, you might want to remove them in case the element is removed from the DOM.

Alternatively, you can also use any JavaScript code you want if the `userProfile.js` is not enough for your needs:

```
document.querySelectorAll(`[data-kcMyCustomValidation]`).forEach((element) => {
    var listener = function (evt) {
        // do something on keyup
    };

    element.addEventListener("keyup", listener);
  });
```

#### [](#managing-attribute-groups)Managing Attribute Groups

At the `Attribute Groups` sub-tab you can create, edit, and delete attribute groups. An attribute group allows you to define a container for correlated attributes so that they are rendered together when at the user-facing forms.

Attribute Group List

![user profile attribute group list](./images/user-profile-attribute-group-list.png)

You can’t delete attribute groups that are bound to attributes. For that, you should first update the attributes to remove the binding.

To create a new group, click on the **Create attributes group** button on the top of the attribute groups listing.

Attribute Group Configuration

![user profile create attribute group](./images/user-profile-create-attribute-group.png)

When configuring the group you can define the following settings:

Name

The name of the attribute, used to uniquely identify an attribute.

Display name

A user-friendly name for the attribute, mainly used when rendering user-facing forms. It also supports [Using Internationalized Messages](#_using-internationalized-messages).

Display description

A user-friendly text that will be displayed as a tooltip when rendering user-facing forms. It also supports [Using Internationalized Messages](#_using-internationalized-messages).

Annotation

In this section, you can associate annotations to the attribute. Annotations are mainly useful to pass over additional metadata to frontends for rendering purposes.

#### [](#_user-profile-json-configuration)Using the JSON configuration

The user profile configuration is stored using a well-defined JSON schema. You can choose from editing the user profile configuration directly by clicking on the `JSON Editor` sub-tab.

JSON Configuration

![user profile json config](./images/user-profile-json-config.png)

The JSON schema is defined as follows:

```
{
  "unmanagedAttributePolicy": "DISABLED",
  "attributes": [
    {
      "name": "myattribute",
      "multivalued": false,
      "displayName": "My Attribute",
      "group": "personalInfo",
      "required": {
        "roles": [ "user", "admin" ],
        "scopes": [ "foo", "bar" ]
      },
      "permissions": {
        "view": [ "admin", "user" ],
        "edit": [ "admin", "user" ]
      },
      "validations": {
        "email": {
          "max-local-length": 64
        },
        "length": {
          "max": 255
        }
      },
      "annotations": {
        "myannotation": "myannotation-value"
      }
    }
  ],
  "groups": [
    {
      "name": "personalInfo",
      "displayHeader": "Personal Information",
      "annotations": {
        "foo": ["foo-value"],
        "bar": ["bar-value"]
      }
    }
  ]
}
```

The schema supports as many attributes and groups as you need.

The `unmanagedAttributePolicy` property defines the unmanaged attribute policy by setting one of following values. For more details, look at the [Understanding Managed and Unmanaged Attributes](#_understanding-managed-and-unmanaged-attributes).

- `DISABLED`
- `ENABLED`
- `ADMIN_VIEW`
- `ADMIN_EDIT`

##### [](#attribute-schema)Attribute Schema

For each attribute you should define a `name` and, optionally, the `required`, `permission`, and the `annotations` settings.

The `required` property defines whether an attribute is required. Keycloak allows you to set an attribute as required based on different conditions.

When the `required` property is defined as an empty object, the attribute is always required.

```
{
  "attributes": [
    {
      "name": "myattribute",
      "required": {}
  ]
}
```

On the other hand, you can choose to make the attribute required only for users, or administrators, or both. As well as mark the attribute as required only in case a specific scope is requested when the user is authenticating in Keycloak.

To mark an attribute as required for a user and/or administrator, set the `roles` property as follows:

```
{
  "attributes": [
    {
      "name": "myattribute",
      "required": {
        "roles": ["user"]
      }
  ]
}
```

The `roles` property expects an array whose values can be either `user` or `admin`, depending on whether the attribute is required by the user or the administrator, respectively.

Similarly, you can choose to make the attribute required when a set of one or more scopes is requested by a client when authenticating a user. For that, you can use the `scopes` property as follows:

```
{
  "attributes": [
    {
      "name": "myattribute",
      "required": {
        "scopes": ["foo"]
      }
  ]
}
```

The `scopes` property is an array whose values can be any string representing a client scope.

The attribute-level `permissions` property can be used to define the read and write permissions to an attribute. The permissions are set based on whether these operations can be performed on the attribute by a user, or administrator, or both.

```
{
  "attributes": [
    {
      "name": "myattribute",
      "permissions": {
        "view": ["admin"],
        "edit": ["user"]
      }
  ]
}
```

Both `view` and `edit` properties expect an array whose values can be either `user` or `admin`, depending on whether the attribute is viewable or editable by the user or the administrator, respectively.

When the `edit` permission is granted, the `view` permission is implicitly granted.

The attribute-level `annotation` property can be used to associate additional metadata to attributes. Annotations are mainly useful for passing over additional information about attributes to frontends rendering user attributes based on the user profile configuration. Each annotation is a key/value pair.

```
{
  "attributes": [
    {
      "name": "myattribute",
      "annotations": {
        "foo": ["foo-value"],
        "bar": ["bar-value"]
      }
  ]
}
```

##### [](#attribute-group-schema)Attribute Group Schema

For each attribute group you should define a `name` and, optionally, the `annotations` settings.

The attribute-level `annotation` property can be used to associate additional metadata to attributes. Annotations are mainly useful for passing over additional information about attributes to frontends rendering user attributes based on the user profile configuration. Each annotation is a key/value pair.

#### [](#customizing-how-uis-are-rendered)Customizing How UIs are Rendered

The UIs from all the user profile contexts (including the administration console) are rendered dynamically accordingly to your user profile configuration.

The default rendering mechanism provides the following capabilities:

- Show or hide fields based on the permissions set to attributes.
- Render markers for required fields based on the constraints set to the attributes.
- Change the field input type (text, date, number, select, multiselect) set to an attribute.
- Mark fields as read-only depending on the permissions set to an attribute.
- Order fields depending on the order set to the attributes.
- Group fields that belong to the same attribute group.
- Dynamically group fields that belong to the same attribute group.

##### [](#ordering-attributes)Ordering attributes

The attribute order is set by dragging and dropping the attribute rows on the attribute listing page.

Ordering Attributes

![user profile attribute list order](./images/user-profile-attribute-list-order.png)

The order you set in this page is respected when fields are rendered in dynamic forms.

##### [](#grouping-attributes)Grouping attributes

When dynamic forms are rendered, they will try to group together attributes that belong to the same attribute group.

Dynamic Update Profile Form

![user profile update profile](./images/user-profile-update-profile.png)

When attributes are linked to an attribute group, the attribute order is also important to make sure attributes within the same group are close together, within a same group header. Otherwise, if attributes within a group do not have a sequential order you might have the same group header rendered multiple times in the dynamic form.

#### [](#enabling-progressive-profiling)Enabling Progressive Profiling

In order to make sure end-user profiles are in compliance with the configuration, administrators can use the `VerifyProfile` required action to eventually force users to update their profiles when authenticating to Keycloak.

The `VerifyProfile` action is similar to the `UpdateProfile` action. However, it leverages all the capabilities provided by the user profile to automatically enforce compliance with the user profile configuration.

When enabled, the `VerifyProfile` action is going to perform the following steps when the user is authenticating:

- Check whether the user profile is fully compliant with the user profile configuration set to the realm. That means running validations and make sure all of them are successful.
- If not, perform an additional step during the authentication so that the user can update any missing or invalid attribute.
- If the user profile is compliant with the configuration, no additional step is performed, and the user continues with the authentication process.

The `VerifyProfile` action is enabled by default. To disable it, click on the `Authentication` link on the left side menu and then click on the `Required Actions` tab. At this tab, use the **Enabled** switch of the `VerifyProfile` action to disable it.

Registering the VerifyProfile Required Action

![user profile register verify profile action](./images/user-profile-register-verify-profile-action.png)

#### [](#_using-internationalized-messages)Using Internationalized Messages

If you want to use internationalized messages when configuring attributes, attributes groups, and annotations, you can set their display name, description, and values, using a placeholder that will translate to a message from a message bundle.

For that, you can use a placeholder to resolve messages keys such as `${myAttributeName}`, where `myAttributeName` is the key for a message in a message bundle. For more details, look at [Localizing messages in a theme](https://www.keycloak.org/ui-customization/localization#_localizing_messages_in_a_theme) about how to add message bundles to custom themes.

### [](#ref-user-credentials_server_administration_guide)Defining user credentials

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/ref-user-credentials.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fref-user-credentials.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fref-user-credentials.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can manage credentials of a user in the **Credentials** tab.

Credential management

![user credentials](./images/user-credentials.png)

You change the priority of credentials by dragging and dropping rows. The new order determines the priority of the credentials for that user. The topmost credential has the highest priority. The priority determines which credential is displayed first after a user logs in.

Type

This column displays the type of credential, for example **password** or **OTP**.

User Label

This is an assignable label to recognize the credential when presented as a selection option during login. It can be set to any value to describe the credential.

Data

This is the non-confidential technical information about the credential. It is hidden, by default. You can click **Show data…​** to display the data for a credential.

Actions

Click **Reset password** to change the password for the user and **Delete** to remove the credential.

You cannot configure other types of credentials for a specific user in the Admin Console; that task is the user’s responsibility.

You can delete the credentials of a user in the event a user loses an OTP device or if credentials have been compromised. You can only delete credentials of a user in the **Credentials** tab.

#### [](#proc-setting-password-user_server_administration_guide)Setting a password for a user

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-setting-password-user.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-setting-password-user.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-setting-password-user.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

If a user does not have a password, or if the password has been deleted, the **Set Password** section is displayed.

If a user already has a password, it can be reset in the **Reset Password** section.

Procedure

1. Click **Users** in the menu. The **Users** page is displayed.
2. Select a user.
3. Click the **Credentials** tab.
4. Type a new password in the **Set Password** section.
5. Click **Set Password**.
   
   If **Temporary** is **ON**, the user must change the password at the first login. To allow users to keep the password supplied, set **Temporary** to **OFF.** The user must click **Set Password** to change the password.

#### [](#requesting-a-user-reset-a-password)Requesting a user reset a password

You can also request that the user reset the password.

Procedure

1. Click **Users** in the menu. The **Users** page is displayed.
2. Select a user.
3. Click the **Credentials** tab.
4. Click **Credential Reset**.
5. Select **Update Password** from the list.
6. Click **Send Email**. The sent email contains a link that directs the user to the **Update Password** window.
7. Optionally, you can set the validity of the email link. This is set to the default preset in the **Tokens** tab in **Realm Settings**.

#### [](#proc_creating-otp_server_administration_guide)Creating an OTP

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-creating-otp.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-creating-otp.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-creating-otp.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

If OTP is conditional in your realm, the user must navigate to Keycloak Account Console to reconfigure a new OTP generator. If OTP is required, then the user must reconfigure a new OTP generator when logging in.

Alternatively, you can send an email to the user that requests the user reset the OTP generator. The following procedure also applies if the user already has an OTP credential.

Prerequisite

- You are logged in to the appropriate realm.

Procedure

1. Click **Users** in the main menu. The **Users** page is displayed.
2. Select a user.
3. Click the **Credentials** tab.
4. Click **Credential Reset**.
5. Set **Reset Actions** to **Configure OTP**.
6. Click **Send Email**. The sent email contains a link that directs the user to the **OTP setup page**.

### [](#con-user-registration_server_administration_guide)Allowing users to self-register

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/con-user-registration.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fcon-user-registration.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fcon-user-registration.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can use Keycloak as a third-party authorization server to manage application users, including users who self-register. If you enable self-registration, the login page displays a registration link so that user can create an account.

Registration link

![registration link](./images/registration-link.png)

A user must add profile information to the registration form to complete registration. The registration form can be customized by removing or adding the fields that must be completed by a user.

Clarification on identity brokering and admin API

Even when self-registrations is disabled, new users can be still added to Keycloak by either:

- Administrator can add new users with the usage of admin console (or admin REST API)
- When identity brokering is enabled, new users authenticated by identity provider may be automatically added/registered in Keycloak storage. See the [First login flow section in the Identity Brokering chapter](#_identity_broker_first_login) for more information.

Also users coming from the [3rd-party user storage](#_user-storage-federation) (for example LDAP) are automatically available in Keycloak when the particular user storage is enabled

Additional resources

- For more information on customizing user registration, see the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/).

#### [](#proc-enabling-user-registration_server_administration_guide)Enabling user registration

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-enabling-user-registration.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-enabling-user-registration.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-enabling-user-registration.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Enable users to self-register.

Procedure

1. Click **Realm Settings** in the main menu.
2. Click the **Login** tab.
3. Toggle **User Registration** to **ON**.

After you enable this setting, a **Register** link displays on the login page of the Admin Console.

#### [](#proc-registering-new-user_server_administration_guide)Registering as a new user

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-registering-new-user.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-registering-new-user.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-registering-new-user.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

As a new user, you must complete a registration form to log in for the first time. You add profile information and a password to register.

Registration form

![registration form](./images/registration-form.png)

Prerequisite

- User registration is enabled.

Procedure

1. Click the **Register** link on the login page. The registration page is displayed.
2. Enter the user profile information.
3. Enter the new password.
4. Click **Register**.

#### [](#proc-requiring-tac-agreement-at-registration_server_administration_guide)Requiring user to agree to terms and conditions during registration

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-requiring-tac-agreement-at-registration.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-requiring-tac-agreement-at-registration.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-requiring-tac-agreement-at-registration.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

For a user to register, you can require agreement to your terms and conditions.

Registration form with required terms and conditions agreement

![registration form with required tac](./images/registration-form-with-required-tac.png)

Prerequisite

- User registration is enabled.
- Terms and conditions required action is enabled.

Procedure

1. Click **Authentication** in the menu. Click the **Flows** tab.
2. Click the **registration** flow.
3. Select **Required** on the **Terms and Conditions** row.
   
   Make the terms and conditions agreement required at registration
   
   ![require tac agreement at registration](./images/require-tac-agreement-at-registration.png)

### [](#con-required-actions_server_administration_guide)Defining actions required at login

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/con-required-actions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fcon-required-actions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fcon-required-actions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can set the actions that a user must perform at the first login. These actions are required after the user provides credentials. After the first login, these actions are no longer required. You add required actions on the **Details** tab of that user.

Some required actions are automatically triggered for the user during login even if they are not explicitly added to this user by the administrator. For example `Update password` action can be triggered if [Password policies](#_password-policies) are configured in a way that the user password needs to be changed every X days. Or `verify profile` action can require the user to update the [User profile](#user-profile) as long as some user attributes do not match the requirements according to the user profile configuration.

The following are examples of required action types:

Update Password

The user must change their password.

Configure OTP

The user must configure a one-time password generator on their mobile device using either the Free OTP or Google Authenticator application.

Verify Email

The user must verify their email account. An email will be sent to the user with a validation link that they must click. Once this workflow is successfully completed, the user will be allowed to log in.

Update Profile

The user must update profile information, such as name, address, email, and phone number.

Some actions do not makes sense to be added to the user account directly. For example, the `Update User Locale` is a helper action to handle some localization related parameters. Another example is the `Delete Credential` action, which is supposed to be triggered as a [Parameterized AIA](#con-aia-parameterized_server_administration_guide). Regarding this one, if the administrator wants to delete the credential of some user, that administrator can do it directly in the Admin Console. The `Delete Credential` action is dedicated to be used for example by the [Keycloak Account Console](#_account-service).

#### [](#proc-setting-required-actions_server_administration_guide)Setting required actions for one user

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-setting-required-actions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-setting-required-actions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-setting-required-actions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can set the actions that are required for any user.

Procedure

1. Click **Users** in the menu.
2. Select a user from the list.
3. Navigate to the **Required User Actions** list.
   
   ![user required action](./images/user-required-action.png)
4. Select all the actions you want to add to the account.
5. Click the **X** next to the action name to remove it.
6. Click **Save** after you select which actions to add.

#### [](#proc-setting-default-required-actions_server_administration_guide)Setting required actions for all users

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-setting-default-required-actions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-setting-default-required-actions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-setting-default-required-actions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can specify what actions are required before the first login of all new users. The requirements apply to a user created by the **Add User** button on the **Users** page or the **Register** link on the login page.

Procedure

1. Click **Authentication** in the menu.
2. Click the **Required Actions** tab.
3. Click the checkbox in the **Set as default action** column for one or more required actions. When a new user logs in for the first time, the selected actions must be executed.

#### [](#proc-enabling-terms-conditions_server_administration_guide)Enabling terms and conditions as a required action

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-enabling-terms-conditions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-enabling-terms-conditions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-enabling-terms-conditions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can enable a required action that new users must accept the terms and conditions before logging in to Keycloak for the first time.

Procedure

1. Click **Authentication** in the menu.
2. Click the **Required Actions** tab.
3. Enable the **Terms and Conditions** action.
4. Edit the `terms.ftl` file in the base login theme.

Additional resources

- For more information on extending and creating themes, see the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/).

### [](#con-aia_server_administration_guide)Application initiated actions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/con-aia.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fcon-aia.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fcon-aia.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Application initiated actions (AIA) allow client applications to request a user to perform an action on the Keycloak side. Usually, when an OIDC client application wants a user to log in, it redirects that user to the login URL as described in the [OIDC section](#con-oidc_server_administration_guide). After login, the user is redirected back to the client application. The user performs the actions that were required by the administrator as described in the [previous section](#proc-setting-required-actions_server_administration_guide) and then is immediately redirected back to the application. However, AIA allows the client application to request some required actions from the user during login. This can be done even if the user is already authenticated on the client and has an active SSO session. It is triggered by adding the `kc_action` parameter to the OIDC login URL with the value containing the requested action. For instance `kc_action=UPDATE_PASSWORD` parameter.

A user may cancel an application initiated action. In this case the user is redirected back to the client application. The redirect URI will contain the query parameters `kc_action_status=cancelled` and `kc_action` with the name of the cancelled action.

The `kc_action` and `kc_action_status` parameters are a Keycloak proprietary mechanism unsupported by the OIDC specification.

Application initiated actions are supported only for OIDC clients.

So if AIA is used, an example flow is similar to the following:

- A client application redirects the user to the OIDC login URL with the additional parameter such as `kc_action=UPDATE_PASSWORD`
- There is a `browser` flow always triggered as described in the [Authentication flows section](#_authentication-flows). If the user was not authenticated, that user needs to authenticate as during normal login. In case the user was already authenticated, that user might be automatically re-authenticated by an SSO cookie without needing to actively re-authenticate and supply the credentials again. In this case, that user will be directly redirected to the screen with the particular action (update password in this case). However, in some cases, active re-authentication is required even if the user has an SSO cookie (See [below](#con-aia-reauth_server_administration_guide) for the details).
- The screen with particular action (in this case `update password`) is displayed to the user, so that user needs to perform a particular action
- Then user is redirected back to the client application

Note that AIA are used by the Keycloak [Account Console](#_account-service) to request update password or to reset other credentials such as OTP or WebAuthn.

Even if the parameter `kc_action` was used, it is not sufficient to assume that the user always performs the action. For example, a user could have manually deleted the `kc_action` parameter from the browser URL. Therefore, no guarantee exists that the user has an OTP for the account after the client requested `kc_action=CONFIGURE_TOTP`. If you want to verify that the user configured two-factor authenticator, the client application may need to check it was configured. For instance by checking the claims like `acr` in the tokens.

#### [](#con-aia-reauth_server_administration_guide)Re-authentication during AIA

In case the user is already authenticated due to an active SSO session, that user usually does not need to actively re-authenticate. However, if that user actively authenticated longer than five minutes ago, the client can still request re-authentication when some AIA is requested. Exceptions exist from this guideline as follows:

- For every required action it is possible to configure the max age on the required action itself in the [Required actions tab](#proc-setting-default-required-actions_server_administration_guide). If the policy is not configured, it defaults to five minutes.
- The action `delete_account` will always require the user to actively re-authenticate
- The action `UPDATE_PASSWORD` might require the user to actively re-authenticate according to the configured [Maximum Authentication Age Password policy](#maximum-authentication-age). In case the policy is not configured, it is also possible to configure it on the required action itself in the [Required actions tab](#proc-setting-default-required-actions_server_administration_guide) when configuring the particular required action. If the policy is not configured in any of those places, it defaults to five minutes.
- If you want to use a shorter re-authentication, you can still use a parameter query parameter such as `max_age` with the specified shorter value or eventually `prompt=login`, which will always require user to actively re-authenticate as described in the OIDC specification. Note that using `max_age` for a longer value than the default five minutes (or the one specifically configured for the required action) is not supported. The `max_age` can be currently used only to make the value shorter than the default five minutes.
- If [Step-up authentication](#_step-up-flow) is enabled and the action is to add or delete a credential, authentication is required with the level corresponding to the given credential. This requirement exists in case the user already has the credential of the particular level. For example, if `otp` and `webauthn` are configured in the authentication flow as 2nd-factor authenticators (both in the authentication flow at level 2) and the user already has a 2nd-factor credential (`otp` or `webauthn` in this case), the user is required to authenticate with an existing 2nd-factor credential to add another 2nd-level credential. In the same manner, deleting an existing 2nd-factor credential (`otp` or `webauthn` in this case), authentication with an existing 2nd-factor level credential is required. The requirement exists for security reasons.

#### [](#con-aia-parameterized_server_administration_guide)Parameterized AIA

Some AIA can require the parameter to be sent together with the action name. For instance, the `Delete Credential` action can be triggered only by AIA and it requires a parameter to be sent together with the name of the action, which points to the ID of the removed credential. So the URL for this example would be `kc_action=delete_credential:ce1008ac-f811-427f-825a-c0b878d1c24b`. In this case, the part after the colon character (`ce1008ac-f811-427f-825a-c0b878d1c24b`) contains the ID of the credential of the particular user, which is to be deleted. The `Delete Credential` action displays the confirmation screen where the user can confirm agreement to delete the credential.

The [Keycloak Account Console](#_account-service) typically uses the `Delete Credential` action when deleting a 2nd-factor credential. You can check the Account Console for examples if you want to use this action directly from your own applications. However, relying on the Account Console is best instead of managing credentials from your own applications.

#### [](#con-aia-available-actions_server_administration_guide)Available actions

To see all available actions, log in to the Admin Console and select `master` realm. Then go to the right top corner and click on the name of the user → select `Realm info` → tab `Provider info`. Then in the table, find SPI `required-action` . In the 2nd column, there are available providers. Those can be used as values of the `kc_action` parameter (unless parameterized as described above). But note that this can be further restricted based on what actions are enabled for your realm in the [Required actions tab](#proc-setting-default-required-actions_server_administration_guide).

### [](#proc-searching-user_server_administration_guide)Searching for a user

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-searching-user.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-searching-user.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-searching-user.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Search for a user to view detailed information about the user, such as the user’s groups and roles.

Prerequisite

- You are in the realm where the user exists.

#### [](#default-search)Default search

Procedure

1. Click **Users** in the main menu. This **Users** page is displayed.
2. Type the full name, last name, first name, or email address of the user you want to search for in the search box. The search returns all users that match your criteria.
   
   The criteria used to match users depends on the syntax used on the search box:
   
   1. `"somevalue"` → performs exact search of the string `"somevalue"`;
   2. `*somevalue*` → performs infix search, akin to a `LIKE '%somevalue%'` DB query;
   3. `somevalue*` or `somevalue` → performs prefix search, akin to a `LIKE 'somevalue%'` DB query.

#### [](#attribute-search)Attribute search

Procedure

1. Click **Users** in the main menu. This **Users** page is displayed.
2. Click **Default search** button and switch it to **Attribute search**.
3. Click **Select attributes** button and specify the attributes to search by.
4. Check **Exact search** checkbox to perform exact match or keep it unchecked to use an infix search for attribute values.
5. Click **Search** button to perform the search. It returns all users that match the criteria.

Searches performed in the **Users** page encompass both Keycloak’s database and configured user federation backends, such as LDAP. Users found in federated backends will be imported into Keycloak’s database if they don’t already exist there.

Additional Resources

- For more information on user federation, see [User Federation](#_user-storage-federation).

### [](#proc-deleting-user_server_administration_guide)Deleting a user

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-deleting-user.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-deleting-user.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-deleting-user.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can delete a user, who no longer needs access to applications. If a user is deleted, the user profile and data is also deleted.

Procedure

1. Click **Users** in the menu. The **Users** page is displayed.
2. Click **View all users** to find a user to delete.
   
   Alternatively, you can use the search bar to find a user.
3. Click **Delete** from the action menu next to the user you want to remove and confirm deletion.

### [](#proc-allow-user-to-delete-account_server_administration_guide)Enabling account deletion by users

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-allow-user-to-delete-account.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-allow-user-to-delete-account.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-allow-user-to-delete-account.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

End users and applications can delete their accounts in the Account Console if you enable this capability in the Admin Console. Once you enable this capability, you can give that capability to specific users.

#### [](#enabling-the-delete-account-capability)Enabling the Delete Account Capability

You enable this capability on the **Required Actions** tab.

Procedure

1. Click **Authentication** in the menu.
2. Click the **Required Actions** tab.
3. Select **Enabled** on the **Delete Account** row.
   
   Delete account on required actions tab
   
   ![enable delete account action](./images/enable-delete-account-action.png)

#### [](#giving-a-user-the-delete-account-role)Giving a user the **delete-account** role

You can give specific users a role that allows account deletion.

Procedure

1. Click **Users** in the menu.
2. Select a user.
3. Click the **Role Mappings** tab.
4. Click the **Assign role** button.
5. Click **account delete-account**.
6. Click **Assign**.
   
   Delete-account role
   
   ![delete-account role](./images/delete-account-client-role.png)

#### [](#deleting-your-account)Deleting your account

Once you have the **delete-account** role, you can delete your own account.

1. Log into the Account Console.
2. At the bottom of the **Personal Info** page, click **Delete Account**.
   
   Delete account page
   
   ![Delete account page](./images/delete-account-page.png)
3. Enter your credentials and confirm the deletion.
   
   Delete confirmation
   
   ![delete account confirm](./images/delete-account-confirm.png)
   
   This action is irreversible. All your data in Keycloak will be removed.

### [](#con-user-impersonation_server_administration_guide)Impersonating a user

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/con-user-impersonation.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fcon-user-impersonation.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fcon-user-impersonation.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

An administrator with the appropriate permissions can impersonate a user. For example, if a user experiences a bug in an application, an administrator can impersonate the user to investigate or duplicate the issue.

Any user with the `impersonation` role in the realm can impersonate a user.

Procedure

1. Click **Users** in the menu.
2. Click a user to impersonate.
3. From the **Actions** list, select **Impersonate**.
   
   ![user impersonate action](./images/user-impersonate-action.png)
   
   - If the administrator and the user are in the same realm, then the administrator will be logged out and automatically logged in as the user being impersonated.
   - If the administrator and user are in different realms, the administrator will remain logged in, and additionally will be logged in as the user in that user’s realm.

In both instances, the **Account Console** of the impersonated user is displayed.

Additional resources

- For more information on assigning administration permissions, see the [Admin Console Access Control](#_admin_permissions) chapter.

### [](#proc-enabling-recaptcha_server_administration_guide)Enabling reCAPTCHA

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/proc-enabling-recaptcha.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fproc-enabling-recaptcha.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fproc-enabling-recaptcha.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To safeguard registration against bots, Keycloak has integration with Google reCAPTCHA (see [Setting up Google reCAPTCHA](#procedure_recaptcha)) and reCAPTCHA Enterprise (see [Setting up Google reCAPTCHA Enterprise](#procedure_recaptcha_enterprise)). The default theme (`register.ftl`) supports both v2 (visible, checkbox-based) and v3 (score-based, invisible) reCAPTCHA (see [Choose the appropriate reCAPTCHA key type](https://docs.cloud.google.com/recaptcha/docs/choose-key-type)).

#### [](#procedure_recaptcha)Setting up Google reCAPTCHA

1. Enter the following URL in a browser:
   
   ```
   https://www.google.com/recaptcha/admin/create
   ```
2. Create a reCAPTCHA and choose between Challenge v2 (visible checkbox) or Score-based, v3 (invisible) to get your reCAPTCHA site key and secret. Note them down for future use in this procedure.
   
   localhost domains are not supported by default. If you wish to continue supporting them for development you can add them to the list of supported domains for your site key.
3. Navigate to the Keycloak admin console.
4. Click **Authentication** in the menu.
5. Click the **Flows** tab.
6. Select **Registration** from the list.
7. Set the **reCAPTCHA** requirement to **Required**. This enables reCAPTCHA.
8. Click the **gear icon** ⚙️ on the **reCAPTCHA** row.
   
   reCAPTCHA config
   
   ![recaptcha config](./images/recaptcha-config.png)
   
   1. Enter the **reCAPTCHA Site Key** generated from the Google reCAPTCHA website.
   2. Enter the **reCAPTCHA Secret** generated from the Google reCAPTCHA website.
   3. Toggle **reCAPTCHA v3** according to your Site Key type: on for score-based reCAPTCHA (v3), off for challenge reCAPTCHA (v2).
   4. (Optional) Toggle **Use recaptcha.net** to use `www.recaptcha.net` instead of `www.google.com` domain for cookies. See [reCAPTCHA faq](https://developers.google.com/recaptcha/docs/faq) for more information.
9. Authorize Google to use the registration page as an iframe.
   
   In Keycloak, websites cannot include a login page dialog in an iframe. This restriction is to prevent clickjacking attacks. You need to change the default HTTP response headers that is set in Keycloak.
   
   1. Click **Realm Settings** in the menu.
   2. Click the **Security Defenses** tab.
   3. Enter `https://www.google.com` in the field for the **X-Frame-Options** header (or `https//www.recaptcha.net` if you enabled **Use recaptcha.net**).
   4. Enter `https://www.google.com` in the field for the **Content-Security-Policy** header (or `https//www.recaptcha.net` if you enabled **Use recaptcha.net**).

#### [](#procedure_recaptcha_enterprise)Setting up Google reCAPTCHA Enterprise

01. Enter the following URL in a browser:
    
    ```
    https://developers.google.com/recaptcha/
    ```
02. Create a key for a "Website" platform, and choose the desired key type. Leave the defaults for v3 reCAPTCHA (invisible), or toggle **Use checkbox challenge** for a v2 reCAPTCHA (visible). Note the site key for future use in this procedure.
    
    The localhost works by default. You do not have to specify a domain.
03. On your Google Cloud Project, go to **Credentials** and create an API key.
    
    For better security, click on **edit api key** and add an API restriction to restrict the key to the **reCAPTCHA Enterprise API** only.
04. Navigate to the Keycloak Admin Console.
05. Click **Authentication** in the menu.
06. Click the **Flows** tab.
07. Duplicate the "registration" flow.
08. Bind the new flow to the **Registration flow**.
09. Edit the new flow:
    
    1. Delete the **reCAPTCHA** step.
    2. Add the step **reCAPTCHA Enterprise** as a sub-step of "registration form" (first step of the flow).
10. Set the **reCAPTCHA Enterprise** requirement to **Required**.
11. Click the **gear icon** ⚙️ on the **reCAPTCHA Enterprise** row.
    
    reCAPTCHA Enterprise config
    
    ![recaptcha enterprise config](./images/recaptcha-enterprise-config.png)
    
    1. Enter the **Recaptcha Project ID** of your Google Cloud console project.
    2. Enter the **Recaptcha Site Key** generated at the beginning of the procedure.
    3. Enter the **Recaptcha API Key** generated at the beginning of the procedure.
    4. Toggle **reCAPTCHA v3** according to your Site Key type: on for score-based reCAPTCHA (v3), off for challenge reCAPTCHA (v2).
    5. (Optional) Customize the **Min. Score Threshold** as you see fit. Set it to the minimum score, between 0.0 and 1.0, that a user should achieve on reCAPTCHA to be allowed to register. See [interpret scores](https://docs.cloud.google.com/recaptcha/docs/interpret-assessment-website#interpret_scores).
    6. (Optional) Toggle **Use recaptcha.net** to use `www.recaptcha.net` instead of `www.google.com` domain for cookies. See [reCAPTCHA faq](https://developers.google.com/recaptcha/docs/faq) for more information.
12. Authorize Google to use the registration page as an iframe. See the last steps of [Setting up Google reCAPTCHA](#procedure_recaptcha) for a detailed procedure.

Additional resources

- For more information on extending and creating themes, see the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/).

### [](#ref-personal-data-collected_server_administration_guide)Personal data collected by Keycloak

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/users/ref-personal-data-collected.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fusers%2Fref-personal-data-collected.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fusers%2Fref-personal-data-collected.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

By default, Keycloak collects the following data:

- Basic user profile data, such as the user email, first name, and last name.
- Basic user profile data used for social accounts and references to the social account when using a social login.
- Device information collected for audit and security purposes, such as the IP address, operating system name, and the browser name.

The information collected in Keycloak is highly customizable. The following guidelines apply when making customizations:

- Registration and account forms can contain custom fields, such as birthday, gender, and nationality. An administrator can configure Keycloak to retrieve data from a social provider or a user storage provider such as LDAP.
- Keycloak collects user credentials, such as password, OTP codes, and WebAuthn public keys. This information is encrypted and saved in a database, so it is not visible to Keycloak administrators. Each type of credential can include non-confidential metadata that is visible to administrators such as the algorithm that is used to hash the password and the number of hash iterations used to hash the password.
- With authorization services and UMA support enabled, Keycloak can hold information about some objects for which a particular user is the owner.

## [](#managing-user-sessions)Managing user sessions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sessions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsessions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsessions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When users log into realms, Keycloak maintains a user session for each user and remembers each client visited by the user within the session. Realm administrators can perform multiple actions on each user session:

- View login statistics for the realm.
- View active users and where they logged in.
- Log a user out of their session.
- Revoke tokens.
- Set up token timeouts.
- Set up session timeouts.

### [](#administering-sessions)Administering sessions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sessions/administering.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsessions%2Fadministering.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsessions%2Fadministering.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To see a top-level view of the active clients and sessions in Keycloak, click **Sessions** from the menu.

Sessions

![Sessions tab](./images/sessions.png)

#### [](#signing-out-all-active-sessions)Signing out all active sessions

You can sign out all users in the realm. From the **Action** list, select **Sign out all active sessions**. All SSO cookies become invalid. Keycloak notifies clients by using the Keycloak OIDC client adapter of the logout event. Clients requesting authentication within active browser sessions must log in again. Client types such as SAML do not receive a back-channel logout request.

Clicking **Sign out all active sessions** does not revoke outstanding access tokens. Outstanding tokens must expire naturally. For clients using the Keycloak OIDC client adapter, you can push a [revocation policy](#_revocation-policy) to revoke the token, but this does not work for other adapters.

#### [](#viewing-client-sessions)Viewing client sessions

Procedure

1. Click **Clients** in the menu.
2. Click a client to see that client’s sessions.
3. Click the **Sessions** tab.
   
   Client sessions
   
   ![Client sessions](./images/client-sessions.png)

#### [](#viewing-user-sessions)Viewing user sessions

Procedure

1. Click **Users** in the menu.
2. Click a user to see that user’s sessions.
3. Click the **Sessions** tab.
   
   User sessions
   
   ![User sessions](./images/user-sessions.png)

### [](#_revocation-policy)Revoking active sessions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sessions/revocation.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsessions%2Frevocation.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsessions%2Frevocation.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

If your system is compromised, you can revoke all active sessions and access tokens.

Procedure

1. Click **Sessions** in the menu.
2. From the **Actions** list, select **Revocation**.
   
   Revocation
   
   ![Revocation](./images/revocation.png)
3. Specify a time and date where sessions or tokens issued before that time and date are invalid using this console.
   
   - Click **Set to now** to set the policy to the current time and date.
   - Click **Push** to push this revocation policy to any registered OIDC client with the Keycloak OIDC client adapter.

### [](#_timeouts)Session and token timeouts

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sessions/timeouts.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsessions%2Ftimeouts.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsessions%2Ftimeouts.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak includes control of the session, cookie, and token timeouts through the **Sessions** and **Tokens** tabs in the **Realm settings** menu.

Sessions tab

![Sessions Tab](./images/sessions-tab.png)

  Configuration Description

SSO Session Idle

This setting is for OIDC clients only. If a user is inactive for longer than this timeout, the user session is invalidated. This timeout value resets when clients request authentication or send a refresh token request. Keycloak adds a window of time to the idle timeout before the session invalidation takes effect. See the [note](#_idle_timeouts_note) later in this section.

SSO Session Max

The maximum time before a user session expires.

SSO Session Idle Remember Me

This setting is only available when **Remember me** is enabled and is similar to the standard SSO Session Idle configuration. Users can specify longer session idle timeouts when they click **Remember Me** when logging in. This setting is an optional configuration and, if its value is not greater than zero, it uses the same idle timeout as the SSO Session Idle configuration.

SSO Session Max Remember Me

This setting is only available when **Remember me** is enabled and is similar to the standard SSO Session Max. Users can specify longer sessions when they click **Remember Me** when logging in. This setting is an optional configuration and, if its value is not greater than zero, it uses the same session lifespan as the SSO Session Max configuration.

Client Session Idle

Idle timeout for the client session. If the user is inactive for longer than this timeout, the client session is invalidated and the refresh token requests bump the idle timeout. This setting never affects the general SSO user session, which is unique. Note the SSO user session is the parent of zero or more client sessions, one client session is created for every different client app the user logs in. This value should specify a shorter idle timeout than the **SSO Session Idle**. Users can override it for individual clients in the **Advanced Settings** client tab. This setting is an optional configuration and, when set to zero, uses the same idle timeout in the SSO Session Idle configuration.

Client Session Max

The maximum time for a client session and before a refresh token expires and invalidates. As in the previous option, this setting never affects the SSO user session and should specify a shorter value than the **SSO Session Max**. Users can override it for individual clients in the **Advanced Settings** client tab. This setting is an optional configuration and, when set to zero, uses the same max timeout in the SSO Session Max configuration.

[]()

Offline Session Idle

This setting is for [offline access](#_offline-access). The amount of time the session remains idle before Keycloak revokes its offline token. Keycloak adds a window of time to the idle timeout before the session invalidation takes effect. See the [note](#_idle_timeouts_note) later in this section.

[]()

Offline Session Max Limited

This setting is for [offline access](#_offline-access). If this flag is **Enabled**, Offline Session Max can control the maximum time the offline token remains active, regardless of user activity. If the flag is **Disabled**, offline sessions never expire by lifespan, only by idle. Once this option is activated, the [Offline Session Max](#_offline-session-max) and [Client Offline Session Max](#_client_offline-session-max) (global option at realm level) can be configured.

[]()

Offline Session Max

This setting is for [offline access](#_offline-access), and it is the maximum time before Keycloak revokes the corresponding offline token. This option controls the maximum amount of time the offline token remains active, regardless of user activity.

[]()

Client Offline Session Max

This setting is for [offline access](#_offline-access), and it is the maximum time before Keycloak revokes the corresponding offline token for the client. This option controls the maximum amount of time the offline token remains active, regardless of user activity. Users can override it for individual clients in the **Advanced Settings** client tab.

Login timeout

The total time a logging in must take. If authentication takes longer than this time, the user must start the authentication process again.

Login action timeout

The Maximum time users can spend on any one page during the authentication process.

Tokens tab

![Tokens Tab](./images/tokens-tab.png)

  Configuration Description

Default Signature Algorithm

The default algorithm used to assign tokens for the realm.

[]()

Revoke Refresh Token

When **Enabled**, Keycloak revokes refresh tokens and issues another token that the client must use. This action applies to OIDC clients performing the refresh token flow.

Access Token Lifespan

When Keycloak creates an OIDC access token, this value controls the lifetime of the token.

Access Token Lifespan For Implicit Flow

With the Implicit Flow, Keycloak does not provide a refresh token. A separate timeout exists for access tokens created by the Implicit Flow.

Client login timeout

The maximum time before clients must finish the Authorization Code Flow in OIDC.

User-Initiated Action Lifespan

The maximum time before a user’s action permission expires. Keep this value short because users generally react to self-created actions quickly.

Default Admin-Initiated Action Lifespan

The maximum time before an action permission sent to a user by an administrator expires. Keep this value long to allow administrators to send e-mails to offline users. An administrator can override the default timeout before issuing the token.

Email Verification

Specifies independent timeout for email verification.

IdP account email verification

Specifies independent timeout for IdP account email verification.

Forgot password

Specifies independent timeout for forgot password.

Execute actions

Specifies independent timeout for execute actions.

The following logic is only applied if persistent user sessions are not active:

For idle timeouts, a two-minute window of time exists that the session is active. For example, when you have the timeout set to 30 minutes, it will be 32 minutes before the session expires.

This action is necessary for some scenarios in cluster and cross-data center environments where the token refreshes on one cluster node a short time before the expiration and the other cluster nodes incorrectly consider the session as expired because they have not yet received the message about a successful refresh from the refreshing node.

### [](#_offline-access)Offline access

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sessions/offline.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsessions%2Foffline.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsessions%2Foffline.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

During [offline access](https://openid.net/specs/openid-connect-core-1_0.html#OfflineAccess) logins, the client application requests an offline token instead of a refresh token. The client application saves this offline token and can use it for future logins if the user logs out. This action is useful if your application needs to perform offline actions on behalf of the user even when the user is not online. For example, a regular data backup.

The client application is responsible for persisting the offline token in storage and then using it to retrieve new access tokens from the Keycloak server.

The difference between a refresh token and an offline token is that an offline token never expires and is not subject to the `SSO Session Idle` timeout and `SSO Session Max` lifespan. The offline token is valid after a user logout. You must use the offline token for a refresh token action at least once per thirty days or for the value of the [Offline Session Idle](#_offline-session-idle).

If you enable [Offline Session Max Limited](#_offline-session-max-limited), offline tokens expire after 60 days even if you use the offline token for a refresh token action. You can change this value, [Offline Session Max](#_offline-session-max), in the Admin Console.

When using offline access, client idle and max timeouts can be overridden at the [client level](#_client_advanced_settings_oidc). The options **Client Offline Session Idle** and **Client Offline Session Max**, in the client **Advanced Settings** tab, allow you to have a shorter offline timeouts for a specific application. Note that client session values also control the refresh token expiration but they never affect the global offline user SSO session. The option **Client Offline Session Max** is only evaluated in the client if [Offline Session Max Limited](#_offline-session-max-limited) is **Enabled** at the realm level.

If you enable the [Revoke Refresh Token](#_revoke-refresh-token) option, you can use each offline token once only. After refresh, you must store the new offline token from the refresh response instead of the previous one.

Users can view and revoke offline tokens that Keycloak grants them in the [User Account Console](#_account-service). Administrators can revoke offline tokens for individual users in the Admin Console in the `Consents` tab. Administrators can view all offline tokens issued in the `Offline Access` tab of each client. Administrators can revoke offline tokens by setting a [revocation policy](#_revocation-policy).

To issue an offline token, users must have the role mapping for the realm-level `offline_access` role. Clients must also have that role in their scope. Clients must add an `offline_access` client scope as an `Optional client scope` to the role, which is done by default.

Clients can request an offline token by adding the parameter `scope=offline_access` when sending their authorization request to Keycloak. The Keycloak OIDC client adapter automatically adds this parameter when you use it to access your application’s secured URL (such as, http://localhost:8080/customer-portal/secured?scope=offline\_access). The Direct Access Grant and Service Accounts support offline tokens if you include `scope=offline_access` in the authentication request body.

Keycloak will limit its internal cache for offline user and offline client sessions to 10000 entries by default, which will reduce the overall memory usage for offline sessions. Items which are evicted from memory will be loaded on-demand from the database when needed. See the server configuration guide to change this default.

### [](#_transient-session)Transient sessions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sessions/transient.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsessions%2Ftransient.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsessions%2Ftransient.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can conduct transient sessions in Keycloak. When using transient sessions, Keycloak does not create a user session after successful authentication. Keycloak creates a temporary, transient session for the scope of the current request that successfully authenticates the user. Keycloak can run [protocol mappers](#_protocol-mappers) using transient sessions after authentication.

The `sid` and `session_state` of the tokens are usually empty when the token is issued with transient sessions. So during transient sessions, the client application cannot refresh tokens or validate a specific session. Sometimes these actions are unnecessary, so you can avoid the additional resource use of persisting user sessions. This session saves performance, memory, and network communication (in cluster and cross-data center environments) resources.

At this moment, transient sessions are automatically used just during [service account authentication](#_service_accounts) with disabled token refresh. Note that token refresh is automatically disabled during service account authentication unless explicitly enabled by client switch `Use refresh tokens for client credentials grant`.

## [](#assigning-permissions-using-roles-and-groups)Assigning permissions using roles and groups

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/assembly-roles-groups.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fassembly-roles-groups.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fassembly-roles-groups.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Roles and groups have a similar purpose, which is to give users access and permissions to use applications. Groups are a collection of users to which you apply roles and attributes. Roles define specific applications permissions and access control.

A role typically applies to one type of user. For example, an organization may include `admin`, `user`, `manager`, and `employee` roles. An application can assign access and permissions to a role and then assign multiple users to that role so the users have the same access and permissions. For example, the Admin Console has roles that give permission to users to access different parts of the Admin Console.

There is a global namespace for roles and each client also has its own dedicated namespace where roles can be defined.

### [](#proc-creating-realm-roles_server_administration_guide)Creating a realm role

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/proc-creating-realm-roles.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-creating-realm-roles.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-creating-realm-roles.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Realm-level roles are a namespace for defining your roles. To see the list of roles, click **Realm Roles** in the menu.

![roles](./images/roles.png)

Procedure

1. Click **Create Role**.
2. Enter a **Role Name**.
3. Enter a **Description**.
4. Click **Save**.

The **description** field can be localized by specifying a substitution variable with `${var-name}` strings. The localized value is configured to your theme within the themes property files. See the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more details.

### [](#con-client-roles_server_administration_guide)Client roles

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/con-client-roles.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fcon-client-roles.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fcon-client-roles.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Client roles are namespaces dedicated to clients. Each client gets its own namespace. Client roles are managed under the **Roles** tab for each client. You interact with this UI the same way you do for realm-level roles.

### [](#_composite-roles)Converting a role to a composite role

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/proc-converting-composite-roles.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-converting-composite-roles.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-converting-composite-roles.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Any realm or client level role can become a *composite role*. A *composite role* is a role that has one or more additional roles associated with it. When a composite role is mapped to a user, the user gains the roles associated with the composite role. This inheritance is recursive so users also inherit any composite of composites. However, we recommend that composite roles are not overused.

Procedure

1. Click **Realm Roles** in the menu.
2. Click the role that you want to convert.
3. From the **Action** list, select **Add associated roles**.

Composite role

![Composite role](./images/composite-role.png)

The role selection UI is displayed on the page and you can associate realm level and client level roles to the composite role you are creating.

In this example, the **employee** realm-level role is associated with the **developer** composite role. Any user with the **developer** role also inherits the **employee** role.

When creating tokens and SAML assertions, any composite also has its associated roles added to the claims and assertions of the authentication response sent back to the client.

### [](#proc-assigning-role-mappings_server_administration_guide)Assigning role mappings

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/proc-assigning-role-mappings.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-assigning-role-mappings.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-assigning-role-mappings.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can assign role mappings to a user through the **Role Mappings** tab for that user.

Procedure

1. Click **Users** in the menu.
2. Click the user that you want to perform a role mapping on.
3. Click the **Role mappings** tab.
4. Click **Assign role**.
5. Select the role you want to assign to the user from the dialog.
6. Click **Assign**.

Role mappings

![Role mappings](./images/user-role-mappings.png)

In the preceding example, we are assigning the composite role **developer** to a user. That role was created in the [Composite Roles](#_composite-roles) topic.

Effective role mappings

![Effective role mappings](./images/effective-role-mappings.png)

When the **developer** role is assigned, the **employee** role associated with the **developer** composite is displayed with **Inherited** "True". **Inherited** roles are the roles explicitly assigned to users and roles that are inherited from composites.

### [](#_default_roles)Using default roles

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/proc-using-default-roles.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-using-default-roles.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-using-default-roles.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Use default roles to automatically assign user role mappings when a user is created or imported through [Identity Brokering](#_identity_broker).

Procedure

1. Click **Realm settings** in the menu.
2. Click the **User registration** tab.
   
   Default roles
   
   ![Default roles](./images/default-roles.png)

This screenshot shows that some *default roles* already exist.

### [](#_role_scope_mappings)Role scope mappings

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/con-role-scope-mappings.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fcon-role-scope-mappings.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fcon-role-scope-mappings.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

On creation of an OIDC access token or SAML assertion, the user role mappings become claims within the token or assertion. Applications use these claims to make access decisions on the resources controlled by the application. Keycloak digitally signs access tokens and applications reuse them to invoke remotely secured REST services. However, these tokens have an associated risk. An attacker can obtain these tokens and use their permissions to compromise your networks. To prevent this situation, use *Role Scope Mappings*.

*Role Scope Mappings* limit the roles declared inside an access token. When a client requests user authentication, the access token it receives contains only the role mappings that are explicitly specified for the client’s scope. The result is that the permissions of each individual access token are limited instead of giving the client access to all the user’s permissions.

By default, each client gets all the role mappings of the user. You can view the role mappings for a client.

Procedure

1. Click **Clients** in the menu.
2. Click the client to go to the details.
3. Click the **Client scopes** tab.
4. Click the link in the row with *Dedicated scope and mappers for this client*
5. Click the **Scope** tab.

Full scope

![Full scope](./images/full-client-scope.png)

By default, the effective roles of scopes are every declared role in the realm. To change this default behavior, toggle **Full Scope Allowed** to **OFF** and declare the specific roles you want in each client. You can also use [client scopes](#_client_scopes) to define the same role scope mappings for a set of clients.

Partial scope

![Partial scope](./images/client-scope.png)

See the [Token Role mappings section](#_oidc_token_role_mappings) for details about the algorithm that adds the roles to the token.

### [](#proc-managing-groups_server_administration_guide)Groups

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/proc-managing-groups.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-managing-groups.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-managing-groups.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Groups in Keycloak manage a common set of attributes and role mappings for each user. Users can be members of any number of groups and inherit the attributes and role mappings assigned to each group.

To manage groups, click **Groups** in the menu.

Groups

![groups](./images/groups.png)

Groups are hierarchical. A group can have multiple subgroups but a group can have only one parent. Subgroups inherit the attributes and role mappings from their parent. Users inherit the attributes and role mappings from their parent as well.

If you have a parent group and a child group, and a user that belongs only to the child group, the user in the child group inherits the attributes and role mappings of both the parent group and the child group.

The hierarchy of a group is sometimes represented using the group path. The path is the complete list of names that represents the hierarchy of a specific group, from top to bottom and separated by slashes `/` (similar to files in a File System). For example a path can be `/top/level1/level2` which means that `top` is a top level group and is parent of `level1`, which in turn is parent of `level2`. This path represents unambiguously the hierarchy for the group `level2`.

Because of historical reasons Keycloak, does not escape slashes in the group name itself. Therefore a group named `level1/group` under `top` uses the path `/top/level1/group`, which is misleading. Keycloak can be started with the option `--spi-group--jpa--escape-slashes-in-group-path` to `true` and then the slashes in the name are escaped with the character `~`. The escape char marks that the slash is part of the name and has no hierarchical meaning. The previous path example would be `/top/level1~/group` when escaped.

```
bin/kc.[sh|bat] start --spi-group--jpa--escape-slashes-in-group-path=true
```

The following example includes a top-level **Sales** group and a child **North America** subgroup.

To add a group:

1. Click the group.
2. Click **Create group**.
3. Enter a group name.
4. Click **Create**.
5. Click the group name.
   
   The group management page is displayed.
   
   Group
   
   ![group](./images/group.png)

Attributes and role mappings you define are inherited by the groups and users that are members of the group.

To add a user to a group:

1. Click **Users** in the menu.
2. Click the user that you want to perform a role mapping on. If the user is not displayed, click **View all users**.
3. Click **Groups**.
4. Click **Join Group**.
5. Select a group from the dialog.
6. Select a group from the **Available Groups** tree.
7. Click **Join**.
   
   Join group
   
   ![user groups](./images/user-groups.png)

To remove a group from a user:

1. Click **Users** in the menu.
2. Click the user to be removed from the group.
3. Click **Leave** on the group table row.

In this example, the user *jimlincoln* is in the *North America* group. You can see *jimlincoln* displayed under the **Members** tab for the group.

Group membership

![group membership](./images/group-membership.png)

#### [](#con-comparing-groups-roles_server_administration_guide)Groups compared to roles

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/con-comparing-groups-roles.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fcon-comparing-groups-roles.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fcon-comparing-groups-roles.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Groups and roles have some similarities and differences. In Keycloak, groups are a collection of users to which you apply roles and attributes. Roles define types of users, and applications assign permissions and access control to roles.

[Composite Roles](#_composite-roles) are similar to Groups as they provide the same functionality. The difference between them is conceptual. Composite roles apply the permission model to a set of services and applications. Use composite roles to manage applications and services.

Groups focus on collections of users and their roles in an organization. Use groups to manage users.

#### [](#proc-specifying-default-groups_server_administration_guide)Using default groups

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/roles-groups/proc-specifying-default-groups.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-specifying-default-groups.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Froles-groups%2Fproc-specifying-default-groups.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To automatically assign group membership to any users who is created or who is imported through [Identity Brokering](#_identity_broker), you use default groups.

1. Click **Realm settings** in the menu.
2. Click the **User registration** tab.
3. Click the **Default Groups** tab.
   
   Default groups
   
   ![Default groups](./images/default-groups.png)

This screenshot shows that some *default groups* already exist.

## [](#configuring-authentication_server_administration_guide)Configuring authentication

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

This chapter covers several authentication topics. These topics include:

- Enforcing strict password and One Time Password (OTP) policies.
- Managing different credential types.
- Logging in with Kerberos.
- Disabling and enabling built-in credential types.

### [](#_password-policies)Password policies

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/password-policies.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fpassword-policies.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fpassword-policies.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When Keycloak creates a realm, it does not associate password policies with the realm. You can set a simple password with no restrictions on its length, security, or complexity. Simple passwords are unacceptable in production environments. Keycloak has a set of password policies available through the Admin Console.

Procedure

1. Click **Authentication** in the menu.
2. Click the **Policies** tab.
3. Select the policy to add in the **Add policy** drop-down box.
4. Enter a value that applies to the policy chosen.
5. Click **Save**.
   
   Password policy
   
   ![Password Policy](./images/password-policy.png)

After saving the policy, Keycloak enforces the policy for new users.

The new policy will not be effective for existing users. Therefore, make sure that you set the password policy from the beginning of the realm creation or add "Update password" to existing users or use "Expire password" to make sure that users update their passwords in next "N" days, which will actually adjust to new password policies.

#### [](#password-policy-types)Password policy types

##### [](#hashalgorithm)HashAlgorithm

Passwords are not stored in cleartext. Before storage or validation, Keycloak hashes passwords using standard hashing algorithms.

Supported password hashing algorithms are shown in the following table.

  Hashing algorithm Description

`argon2`

Argon2 (default for non-FIPS deployments)

`pbkdf2-sha512`

PBKDF2 with SHA512 (default for FIPS deployments)

`pbkdf2-sha256`

PBKDF2 with SHA256

`pbkdf2`

PBKDF2 with SHA1 (deprecated)

It is highly recommended to use Argon2 when possible as it has significantly less CPU requirements compared to PBKDF2, while at the same time being more secure.

The default password hashing algorithm for the server can be configured with `--spi-password-hashing--provider-default=<algorithm>`.

To prevent excessive memory and CPU usage, the parallel computation of hashes by Argon2 is by default limited to the number of cores available to the JVM. To configure the Argon2 hashing provider, use its provider options.

See the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) on how to add your own hashing algorithm.

If you change the hashing algorithm, password hashes in storage will not change until the user logs in.

##### [](#hashing-iterations)Hashing iterations

Specifies the number of times Keycloak hashes passwords before storage or verification. The default value is -1, which uses the default hashing iterations for the selected hashing algorithm as listed in the following table.

  Hashing algorithm Default hash iterations

`argon2`

5

`pbkdf2-sha512`

210,000

`pbkdf2-sha256`

600,000

`pbkdf2`

1,300,000

In most cases the hashing iterations should not be changed from the recommended default values. Lower values for iterations provide insufficient security, while higher values result in higher CPU power requirements.

##### [](#digits)Digits

The number of numerical digits required in the password string.

##### [](#lowercase-characters)Lowercase characters

The number of lower case letters required in the password string.

##### [](#uppercase-characters)Uppercase characters

The number of upper case letters required in the password string.

##### [](#special-characters)Special characters

The number of special characters required in the password string.

##### [](#not-username)Not username

The password cannot be the same as the username.

##### [](#not-email)Not email

The password cannot be the same as the email address of the user.

##### [](#regular-expression)Regular expression

Password must match one or more defined Java regular expression patterns. See [Java’s regular expression documentation](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/regex/Pattern.html) for the syntax of those expressions.

##### [](#expire-password)Expire password

The number of days the password is valid. When the number of days has expired, the user must change their password.

##### [](#not-recently-used)Not recently used

Password cannot be already used by the user. Keycloak stores a history of used passwords. The number of old passwords stored is configurable in Keycloak.

##### [](#not-recently-used-in-days)Not recently used (In Days)

Password cannot be reused within the configured time period (in days). If the new password was last set within this period, the user will be forced to provide a different one.

##### [](#password-blacklist)Password blacklist

Password must not be in a blacklist file.

- Blacklist files are UTF-8 plain-text files with Unix line endings. Every line represents a blacklisted password.
- Keycloak compares passwords in a case-insensitive manner.
- The value of the blacklist file must be the name of the blacklist file, for example, `100k_passwords.txt`.
- Blacklist files resolve against `${kc.home.dir}/data/password-blacklists/` by default. Customize this path using:
  
  - The `keycloak.password.blacklists.path` system property.
  - The `blacklistsPath` property of the `passwordBlacklist` policy SPI configuration. To configure the blacklist folder using the CLI, use `--spi-password-policy--password-blacklist--blacklists-path=/path/to/blacklistsFolder`.
- Blacklist files are reloaded automatically when modified.
  
  - The file is checked for updates during password policy validation but at most once every 60 seconds. Use `--spi-password-policy--password-blacklist--check-interval-seconds=<seconds>` to change the default (0 to disable).
  - The blacklist file must be updated atomically (write to a temp file, then rename) to avoid reading a partially written file.

A note about False Positives

The current implementation uses a BloomFilter for fast and memory efficient containment checks, such as whether a given password is contained in a blacklist, with the possibility for false positives.

- By default a false positive probability of `0.01%` is used.
- To change the false positive probability by CLI configuration, use `--spi-password-policy--password-blacklist--false-positive-probability=0.00001`.

##### [](#maximum-authentication-age)Maximum Authentication Age

Specifies the maximum age of a user authentication in seconds with which the user can update a password without re-authentication. A value of `0` indicates that the user has to always re-authenticate with their current password before they can update the password. See [AIA section](#con-aia-reauth_server_administration_guide) for some additional details about this policy.

The Maximum Authentication Age is configurable also when configuring the required action **Update Password** in the **Required Actions** tab in the Admin Console. The better choice is to use the required action for the configuration because the *Maximum Authentication Age* password policy might be deprecated/removed in the future.

### [](#one-time-password-otp-policies)One Time Password (OTP) policies

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/otp-policies.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fotp-policies.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fotp-policies.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak has several policies for setting up a FreeOTP or Google Authenticator One-Time Password generator.

Procedure

1. Click **Authentication** in the menu.
2. Click the **Policy** tab.
3. Click the **OTP Policy** tab.

Otp Policy

![OTP Policy](./images/otp-policy.png)

Keycloak generates a QR code on the OTP set-up page, based on information configured in the **OTP Policy** tab. FreeOTP and Google Authenticator scan the QR code when configuring OTP.

#### [](#time-based-or-counter-based-one-time-passwords)Time-based or counter-based one time passwords

The algorithms available in Keycloak for your OTP generators are time-based and counter-based.

With Time-Based One Time Passwords (TOTP), the token generator will hash the current time and a shared secret. The server validates the OTP by comparing the hashes within a window of time to the submitted value. TOTPs are valid for a short window of time.

With Counter-Based One Time Passwords (HOTP), Keycloak uses a shared counter rather than the current time. The Keycloak server increments the counter with each successful OTP login. Valid OTPs change after a successful login.

TOTP is more secure than HOTP because the matchable OTP is valid for a short window of time, while the OTP for HOTP is valid for an indeterminate amount of time. HOTP is more user-friendly than TOTP because no time limit exists to enter the OTP.

HOTP requires a database update every time the server increments the counter. This update is a performance drain on the authentication server during heavy load. To increase efficiency, TOTP does not remember passwords used, so there is no need to perform database updates. The drawback is that it is possible to reuse TOTPs in the valid time interval.

#### [](#totp-configuration-options)TOTP configuration options

##### [](#otp-hash-algorithm)OTP hash algorithm

The default algorithm is SHA1. The other, more secure options are SHA256 and SHA512.

##### [](#number-of-digits)Number of digits

The length of the OTP. Short OTPs are user-friendly, easier to type, and easier to remember. Longer OTPs are more secure than shorter OTPs.

##### [](#look-around-window)Look around window

The number of intervals the server attempts to match the hash. This option is present in Keycloak if the clock of the TOTP generator or authentication server becomes out-of-sync. The default value of 1 is adequate. For example, if the time interval for a token is 30 seconds, the default value of 1 means it will accept valid tokens in the 90-second window (time interval 30 seconds + look ahead 30 seconds + look behind 30 seconds). Every increment of this value increases the valid window by 60 seconds (look ahead 30 seconds + look behind 30 seconds).

##### [](#otp-token-period)OTP token period

The time interval in seconds the server matches a hash. Each time the interval passes, the token generator generates a TOTP.

##### [](#reusable-code)Reusable code

Determine whether OTP tokens can be reused in the authentication process or user needs to wait for the next token. Users cannot reuse those tokens by default, and the administrator needs to explicitly specify that those tokens can be reused.

#### [](#hotp-configuration-options)HOTP configuration options

##### [](#otp-hash-algorithm-2)OTP hash algorithm

The default algorithm is SHA1. The other, more secure options are SHA256 and SHA512.

##### [](#number-of-digits-2)Number of digits

The length of the OTP. Short OTPs are user-friendly, easier to type, and easier to remember. Longer OTPs are more secure than shorter OTPs.

##### [](#look-around-window-2)Look around window

The number of previous and following intervals the server attempts to match the hash. This option is present in Keycloak if the clock of the TOTP generator or authentication server become out-of-sync. The default value of 1 is adequate. This option is present in Keycloak to cover when the user’s counter gets ahead of the server.

##### [](#initial-counter)Initial counter

The value of the initial counter.

### [](#_authentication-flows)Authentication flows

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/flows.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fflows.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fflows.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

An *authentication flow* is a container of authentications, screens, and actions, during log in, registration, and other Keycloak workflows.

#### [](#built-in-flows)Built-in flows

Keycloak has several built-in flows. You cannot modify these flows, but you can alter the flow’s requirements to suit your needs.

Procedure

1. Click **Authentication** in the menu.
2. Click on the *Browser* item in the list to see the details.

Browser flow

![Browser Flow](./images/browser-flow.png)

##### [](#auth-type)Auth type

The name of the authentication or the action to execute. If an authentication is indented, it is in a sub-flow. It may or may not be executed, depending on the behavior of its parent.

1. Cookie
   
   The first time a user logs in successfully, Keycloak sets a session cookie. If the cookie is already set, this authentication type is successful. Since the cookie provider returned success and each execution at this level of the flow is *alternative*, Keycloak does not perform any other execution. This results in a successful login.
2. Kerberos
   
   This authenticator is disabled by default and is skipped during the Browser Flow.
3. Identity Provider Redirector
   
   This action is configured through the **Actions** &gt; **Config** link. It redirects to another IdP for [identity brokering](#_identity_broker).
4. Forms
   
   Since this sub-flow is marked as *alternative*, it will not be executed if the **Cookie** authentication type passed. This sub-flow contains an additional authentication type that needs to be executed. Keycloak loads the executions for this sub-flow and processes them.

The first execution is the **Username Password Form**, an authentication type that renders the username and password page. It is marked as *required*, so the user must enter a valid username and password.

The second execution is the **Browser - Conditional 2FA** sub-flow. This sub-flow is *conditional* and executes depending on the result of the **Condition - User Configured** and **Condition - credential** execution. If the result is true, Keycloak loads the executions for this sub-flow and processes them.

The next execution is the **Condition - User Configured** authentication. This authentication checks if Keycloak has configured other executions in the flow for the user. The **Browser - Conditional 2FA** sub-flow executes only when the user has a configured OTP credential.

The next execution is the **Condition - credential** authentication. The step checks if the authentication process has already authenticated the user with a passwordless WebAuthn credential (passkey), avoiding 2FA in that case.

The final execution is the **OTP Form**. Keycloak marks this execution as *alternative*, but it runs only when the user has an OTP credential set up because of the setup in the *conditional* sub-flow. If the OTP credential is not set up, the user does not see an OTP form.

The default **browser** flow contains two more executions inside the **Browser - Conditional 2FA**, **WebAuthn Authenticator** and **Recovery Authentication Code Form**. These executions are *Disabled* by default and they are the other 2FA methods that can be added to the flow. Change the requirement from *Disabled* to *Alternative* to make them available if the respective credential has been configured for the user. If the user has configured all alternative credential types, the credential with the highest priority is displayed by default. However, the **Try Another Way** option will appear so that the user has the alternative methods to log in.

##### [](#_execution-requirements)Requirement

A drop-down menu that controls the execution of an action.

Required

All *Required* elements in the flow must be successfully sequentially executed. The flow terminates if a required element fails.

Alternative

Only a single element must successfully execute for the flow to evaluate as successful. Because the *Required* flow elements are sufficient to mark a flow as successful, any *Alternative* flow element within a flow containing *Required* flow elements will not execute.

Disabled

The element does not count to mark a flow as successful.

Conditional

This requirement type is only set on sub-flows.

- A *Conditional* sub-flow contains executions. These executions must evaluate to logical statements.
- If all executions evaluate as *true*, the *Conditional* sub-flow acts as *Required*.
- If any executions evaluate as *false*, the *Conditional* sub-flow acts as *Disabled*.
- If you do not set an execution, the *Conditional* sub-flow acts as *Disabled*.
- If a flow contains executions and the flow is not set to *Conditional*, Keycloak does not evaluate the executions, and the executions are considered functionally *Disabled*.

#### [](#creating-flows)Creating flows

Important functionality and security considerations apply when you design a flow.

To create a flow, perform the following:

Procedure

1. Click **Authentication** in the menu.
2. Click **Create flow**.

You can copy and then modify an existing flow. Click the "Action list" (the three dots at the end of the row), click **Duplicate**, and enter a name for the new flow.

When creating a new flow, you must create a top-level flow first with the following options:

Name

The name of the flow.

Description

The description you can set to the flow.

Top-Level Flow Type

The type of flow. The type **client** is used only for the authentication of clients (applications). For all other cases, choose **basic**.

Create a top-level flow

![Top Level Flow](./images/Create-top-level-flow.png)

When Keycloak has created the flow, Keycloak displays the **Add step**, and **Add sub-flow** buttons.

An empty new flow

![New Flow](./images/New-flow.png)

Three factors determine the behavior of flows and sub-flows.

- The structure of the flow and sub-flows.
- The executions within the flows
- The requirements set within the sub-flows and the executions.

Executions have a wide variety of actions, from sending a reset email to validating an OTP. Add executions with the **Add step** button.

Adding an authentication execution

![Adding an Authentication Execution](./images/Create-authentication-execution.png)

Authentication executions can optionally have a reference value configured. This can be utilized by the *Authentication Method Reference (AMR)* protocol mapper to populate the *amr* claim in OIDC access and ID tokens (for more information on the AMR claim, see [RFC-8176](https://www.rfc-editor.org/rfc/rfc8176.html)). When the *Authentication Method Reference (AMR)* protocol mapper is configured for a client, it will populate the *amr* claim with the reference value for any authenticator execution the user successfully completes during the authentication flow.

Adding an authenticator reference value

![Configuring an Authenticator Reference Value](./images/config-authenticator-reference.png)

Two types of executions exist, *automatic executions* and *interactive executions*. *Automatic executions* are similar to the **Cookie** execution and will automatically perform their action in the flow. *Interactive executions* halt the flow to get input. Executions executing successfully set their status to *success*. For a flow to complete, it needs at least one execution with a status of *success*.

You can add sub-flows to top-level flows with the **Add sub-flow** button. The **Add sub-flow** button displays the **Create Execution Flow** page. This page is similar to the **Create Top Level Form** page. The difference is that the **Flow Type** can be **basic** (default) or **form**. The **form** type constructs a sub-flow that generates a form for the user, similar to the built-in **Registration** flow. Sub-flows success depends on how their executions evaluate, including their contained sub-flows. See the [execution requirements section](#_execution-requirements) for an in-depth explanation of how sub-flows work.

After adding an execution, check the requirement has the correct value.

All elements in a flow have a **Delete** option next to the element. Some executions have a **⚙️** menu item (the gear icon) to configure the execution. It is also possible to add executions and sub-flows to sub-flows with the **Add step** and **Add sub-flow** links.

Since the order of execution is important, you can move executions and sub-flows up and down by dragging their names.

Make sure to properly test your configuration when you configure the authentication flow to confirm that no security holes exist in your setup. We recommend that you test various corner cases. For example, consider testing the authentication behavior for a user when you remove various credentials from the user’s account before authentication.

As an example, when 2nd-factor authenticators, such as OTP Form or WebAuthn Authenticator, are configured in the flow as REQUIRED and the user does not have credential of particular type, the user will be able to set up the particular credential during authentication itself. This situation means that the user does not authenticate with this credential as he set up it right during the authentication. So for browser authentication, make sure to configure your authentication flow with some 1st-factor credentials such as Password or WebAuthn Passwordless Authenticator.

#### [](#creating-a-password-less-browser-login-flow)Creating a password-less browser login flow

To illustrate the creation of flows, this section describes creating an advanced browser login flow. The purpose of this flow is to allow a user a choice between logging in using a password-less manner with [WebAuthn](#webauthn_server_administration_guide), or two-factor authentication with a password and OTP.

Procedure

01. Click **Authentication** in the menu.
02. Click the **Flows** tab.
03. Click **Create flow**.
04. Enter `Browser Password-less` as a name.
05. Click **Create**.
06. Click **Add execution**.
07. Select **Cookie** from the list.
08. Click **Add**.
09. Select **Alternative** for the **Cookie** authentication type to set its requirement to alternative.
10. Click **Add step**.
11. Select **Kerberos** from the list.
12. Click **Add**.
13. Click **Add step**.
14. Select **Identity Provider Redirector** from the list.
15. Click **Add**.
16. Select **Alternative** for the **Identity Provider Redirector** authentication type to set its requirement to alternative.
17. Click **Add sub-flow**.
18. Enter **Forms** as a name.
19. Click **Add**.
20. Select **Alternative** for the **Forms** authentication type to set its requirement to alternative.
    
    The common part with the browser flow
    
    ![Passwordless browser login](./images/Passwordless-browser-login-common.png)
21. Click **+** menu of the **Forms** execution.
22. Select **Add step**.
23. Select **Username Form** from the list.
24. Click **Add**.

At this stage, the form requires a username but no password. We must enable password authentication to avoid security risks.

01. Click **+** menu of the **Forms** sub-flow.
02. Click **Add sub-flow**.
03. Enter `Authentication` as name.
04. Click **Add**.
05. Select **Required** for the **Authentication** authentication type to set its requirement to required.
06. Click **+** menu of the **Authentication** sub-flow.
07. Click **Add step**.
08. Select **WebAuthn Passwordless Authenticator** from the list.
09. Click **Add**.
10. Select **Alternative** for the **Webauthn Passwordless Authenticator** authentication type to set its requirement to alternative.
11. Click **+** menu of the **Authentication** sub-flow.
12. Click **Add sub-flow**.
13. Enter `Password with OTP` as name.
14. Click **Add**.
15. Select **Alternative** for the **Password with OTP** authentication type to set its requirement to alternative.
16. Click **+** menu of the **Password with OTP** sub-flow.
17. Click **Add step**.
18. Select **Password Form** from the list.
19. Click **Add**.
20. Select **Required** for the **Password Form** authentication type to set its requirement to required.
21. Click **+** menu of the **Password with OTP** sub-flow.
22. Click **Add step**.
23. Select **OTP Form** from the list.
24. Click **Add**.
25. Click **Required** for the **OTP Form** authentication type to set its requirement to required.

Finally, change the bindings.

1. Click the **Action** menu at the top of the screen.
2. Select **Bind flow** from the menu.
3. Click the **Browser Flow** drop-down list.
4. Click **Save**.

A password-less browser login

![Passwordless browser login](./images/Passwordless-browser-login.png)

After entering the username, the flow works as follows:

If users have WebAuthn passwordless credentials recorded, they can use these credentials to log in directly. This is the password-less login. The user can also select **Password with OTP** because the `WebAuthn Passwordless` execution and the `Password with OTP` flow are set to **Alternative**. If they are set to **Required**, the user has to enter WebAuthn, password, and OTP.

If the user selects the **Try another way** link with `WebAuthn passwordless` authentication, the user can choose between `Password` and `Passkey` (WebAuthn passwordless). When selecting the password, the user will need to continue and log in with the assigned OTP. If the user has no WebAuthn credentials, the user must enter the password and then the OTP. If the user has no OTP credential, they will be asked to record one.

Since the WebAuthn Passwordless execution is set to **Alternative** rather than **Required**, this flow will never ask the user to register a WebAuthn credential. For a user to have a Webauthn credential, an administrator must add a required action to the user. Do this by:

1. Enabling the **Webauthn Register Passwordless** required action in the realm (see the [WebAuthn](#webauthn_server_administration_guide) documentation).
2. Setting the required action using the **Credential Reset** part of a user’s [Credentials](#ref-user-credentials_server_administration_guide) management menu.

Creating an advanced flow such as this can have side effects. For example, if you enable the ability to reset the password for users, this would be accessible from the password form. In the default `Reset Credentials` flow, users must enter their username. Since the user has already entered a username earlier in the `Browser Password-less` flow, this action is unnecessary for Keycloak and suboptimal for user experience. To correct this problem, you can:

- Duplicate the `Reset Credentials` flow. Set its name to `Reset Credentials for password-less`, for example.
- Click **Delete** (trash icon) of the **Choose user** step.
- In the **Action** menu, select **Bind flow** and select **Reset credentials flow** from the dropdown and click **Save**

#### [](#_client-policy-auth-flow)Using Client Policies to Select an Authentication Flow

[Client Policies](#_client_policies) can be used to dynamically select an Authentication Flow based on specific conditions, such as requesting a particular scope or an ACR (Authentication Context Class Reference) using the `AuthenticationFlowSelectorExecutor` in combination with the condition you prefer.

The `AuthenticationFlowSelectorExecutor` allows you to select an appropriate authentication flow and set the level of authentication to be applied once the selected flow is completed.

A possible configuration involves using the `ACRCondition` in combination with the `AuthenticationFlowSelectorExecutor`. This setup enables you to select an authentication flow based on the requested ACR and have the ACR value included in the token using [ACR to LoA Mapping](#_mapping-acr-to-loa-realm).

For more details, see [Client Policies](#_client_policies).

#### [](#_step-up-flow)Creating a browser login flow with step-up mechanism

This section describes how to create advanced browser login flow using the step-up mechanism. The purpose of step-up authentication is to allow access to clients or resources based on a specific authentication level of a user.

Procedure

01. Click **Authentication** in the menu.
02. Click the **Flows** tab.
03. Click **Create flow**.
04. Enter `Browser Incl Step up Mechanism` as a name.
05. Click **Save**.
06. Click **Add execution**.
07. Select **Cookie** from the list.
08. Click **Add**.
09. Select **Alternative** for the **Cookie** authentication type to set its requirement to alternative.
10. Click **Add sub-flow**.
11. Enter **Auth Flow** as a name.
12. Click **Add**.
13. Click **Alternative** for the **Auth Flow** authentication type to set its requirement to alternative.

Now you configure the flow for the first authentication level.

01. Click **+** menu of the **Auth Flow**.
02. Click **Add sub-flow**.
03. Enter `1st Condition Flow` as a name.
04. Click **Add**.
05. Click **Conditional** for the **1st Condition Flow** authentication type to set its requirement to conditional.
06. Click **+** menu of the **1st Condition Flow**.
07. Click **Add condition**.
08. Select **Conditional - Level Of Authentication** from the list.
09. Click **Add**.
10. Click **Required** for the **Conditional - Level Of Authentication** authentication type to set its requirement to required.
11. Click **⚙️** (gear icon).
12. Enter `Level 1` as an alias.
13. Enter `1` for the Level of Authentication (LoA).
14. Set Max Age to **36000**. This value is in seconds and it is equivalent to 10 hours, which is the default `SSO Session Max` timeout set in the realm. As a result, when a user authenticates with this level, subsequent SSO logins can reuse this level and the user does not need to authenticate with this level until the end of the user session, which is 10 hours by default.
15. Click **Save**
    
    Configure the condition for the first authentication level
    
    ![Authentication step up condition 1](./images/authentication-step-up-condition-1.png)
16. Click **+** menu of the **1st Condition Flow**.
17. Click **Add step**.
18. Select **Username Password Form** from the list.
19. Click **Add**.

Now you configure the flow for the second authentication level.

01. Click **+** menu of the **Auth Flow**.
02. Click **Add sub-flow**.
03. Enter `2nd Condition Flow` as an alias.
04. Click **Add**.
05. Click **Conditional** for the **2nd Condition Flow** authentication type to set its requirement to conditional.
06. Click **+** menu of the **2nd Condition Flow**.
07. Click **Add condition**.
08. Select **Conditional - Level Of Authentication** from the item list.
09. Click **Add**.
10. Click **Required** for the **Conditional - Level Of Authentication** authentication type to set its requirement to required.
11. Click **⚙️** (gear icon).
12. Enter `Level 2` as an alias.
13. Enter `2` for the Level of Authentication (LoA).
14. Set Max Age to **0**. As a result, when a user authenticates, this level is valid just for the current authentication, but not any subsequent SSO authentications. So the user will always need to authenticate again with this level when this level is requested.
15. Click **Save**
    
    Configure the condition for the second authentication level
    
    ![Autehtnication step up condition 2](./images/authentication-step-up-condition-2.png)
16. Click **+** menu of the **2nd Condition Flow**.
17. Click **Add step**.
18. Select **OTP Form** from the list.
19. Click **Add**.
20. Click **Required** for the **OTP Form** authentication type to set its requirement to required.

Finally, change the bindings.

1. Click the **Action** menu at the top of the screen.
2. Select **Bind flow** from the list.
3. Select **Browser Flow** in the dropdown.
4. Click **Save**.

Browser login with step-up mechanism

![Authentication step up flow](./images/authentication-step-up-flow.png)

Request a certain authentication level with OpenID Connect

To use the step-up mechanism, you specify a requested level of authentication (LoA) in your authentication request. The `claims` parameter is used for this purpose:

```
https://{DOMAIN}/realms/{REALMNAME}/protocol/openid-connect/auth?client_id={CLIENT-ID}&redirect_uri={REDIRECT-URI}&scope=openid&response_type=code&response_mode=query&nonce=exg16fxdjcu&claims=%7B%22id_token%22%3A%7B%22acr%22%3A%7B%22essential%22%3Atrue%2C%22values%22%3A%5B%22gold%22%5D%7D%7D%7D
```

The `claims` parameter is specified in a JSON representation:

```
claims= {
            "id_token": {
                "acr": {
                    "essential": true,
                    "values": ["gold"]
                }
            }
        }
```

The Keycloak javascript adapter has support for easy construct of this JSON and sending it in the login request. See **Keycloak JavaScript adapter** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section for more details.

You can also use simpler parameter `acr_values` instead of `claims` parameter to request particular levels as non-essential. This is mentioned in the OIDC specification.

You can also configure the default level for the particular client, which is used when the parameter `acr_values` or the parameter `claims` with the `acr` claim is not present. For further details, see [Client ACR configuration](#_mapping-acr-to-loa-client)).

To request the acr\_values as text (such as `gold`) instead of a numeric value, you configure the mapping between the ACR and the LoA. It is possible to configure it at the realm level (recommended) or at the client level. For configuration see [ACR to LoA Mapping](#_mapping-acr-to-loa-realm).

For more details see the [official OIDC specification](https://openid.net/specs/openid-connect-core-1_0.html#acrSemantics).

**Flow logic**

The logic for the previous configured authentication flow is as follows:  
If a client request a high authentication level, meaning Level of Authentication 2 (LoA 2), a user has to perform full 2-factor authentication: Username/Password + OTP. However, if a user already has a session in Keycloak, that was logged in with username and password (LoA 1), the user is only asked for the second authentication factor (OTP).

The option **Max Age** in the condition determines how long (how much seconds) the subsequent authentication level is valid. This setting helps to decide whether the user will be asked to present the authentication factor again during a subsequent authentication. If the particular level X is requested by the `claims` or `acr_values` parameter and user already authenticated with level X, but it is expired (for example max age is configured to 300 and user authenticated before 310 seconds) then the user will be asked to re-authenticate again with the particular level. However if the level is not yet expired, the user will be automatically considered as authenticated with that level.

Using **Max Age** with the value 0 means, that particular level is valid just for this single authentication. Hence every re-authentication requesting that level will need to authenticate again with that level. This is useful for operations that require higher security in the application (e.g. send payment) and always require authentication with the specific level.

Note that parameters such as `claims` or `acr_values` might be changed by the user in the URL when the login request is sent from the client to the Keycloak via the user’s browser. This situation can be mitigated if client uses PAR (Pushed authorization request), a request object, or other mechanisms that prevents the user from rewrite the parameters in the URL. Hence after the authentication, clients are encouraged to check the ID Token to double-check that `acr` in the token corresponds to the expected level.

If no explicit level is requested by parameters, the Keycloak will require the authentication with the first LoA condition found in the authentication flow, such as the Username/Password in the preceding example. When a user was already authenticated with that level and that level expired, the user is not required to re-authenticate, but `acr` in the token will have the value 0. This result is considered as authentication based solely on `long-lived browser cookie` as mentioned in the section 2 of OIDC Core 1.0 specification.

During the first authentication of the user, the first configured subflow with the **Conditional - Level Of Authentication** is always executed (regardless of the requested level) as the user does not yet have any level. Therefore, we recommend that the first level subflow contains the minimal required authenticators for user authentication. In addition, ensure that the subflows with different values of **Conditional - Level Of Authentication** are ordered starting with the lowest as shown in the example above. For example, if you configure a subflow with level 2 and then add another subflow with level 1, the level 2 subflow will be always asked during the first authentication, which may not be the desired behavior.

A conflict situation may arise when an admin specifies several flows, sets different LoA levels to each, and assigns the flows to different clients. However, the rule is always the same: if a user has a certain level, it needs only have that level to connect to a client. It’s up to the admin to make sure that the LoA is coherent.

Step-up authentication with Level of Authentication conditions is intended for use cases where each level requires all authentication methods from the preceding levels. For instance, level X must always include all authentication methods required by level X-1. For use cases where a specific level, such as level 3, requires a different authentication method from the previous levels, it may be more appropriate to use mapping of ACR to a specific flow. For more details, see [Using Client Policies to Select an Authentication Flow](#_client-policy-auth-flow).

**Example scenario**

1. Max Age is configured as 300 seconds for level 1 condition.
2. Login request is sent without requesting any acr. Level 1 will be used and the user needs to authenticate with username and password. The token will have `acr=1`.
3. Another login request is sent after 100 seconds. The user is automatically authenticated due to the SSO and the token will return `acr=1`.
4. Another login request is sent after another 201 seconds (301 seconds since authentication in point 2). The user is automatically authenticated due to the SSO, but the token will return `acr=0` due the level 1 is considered expired.
5. Another login request is sent, but now it will explicitly request ACR of level 1 in the `claims` parameter. User will be asked to re-authenticate with username/password and then `acr=1` will be returned in the token.

**ACR claim in the token**

ACR claim is added to the token by the `acr loa level` protocol mapper defined in the `acr` client scope. This client scope is the realm default client scope and hence will be added to all newly created clients in the realm.

In case you do not want `acr` claim inside tokens or you need some custom logic for adding it, you can remove the client scope from your client.

Note when the login request initiates a request with the `claims` parameter requesting `acr` as `essential` claim, then Keycloak will always return one of the specified levels. If it is not able to return one of the specified levels (For example if the requested level is unknown or bigger than configured conditions in the authentication flow), then Keycloak will throw an error.

#### [](#_step-up-authentication-saml)Step-up authentication for SAML

Step-up Authentication for SAML is **Preview** and is not fully supported. This feature is disabled by default.

To enable start the server with `--features=preview` or `--features=step-up-authentication-saml`

For the SAML protocol, the step-up authentication uses the `<AuthnContextClassRef>` element (Authentication Context Class Reference or ACR) to map the Level of Authentication (LoA). This element is a URI reference that identifies an authentication context declaration. The LoA is requested by the client in the SAML request via the `<RequestedAuthnContext>` element.

```
<samlp:AuthnRequest xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol" xmlns="urn:oasis:names:tc:SAML:2.0:assertion" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion" AssertionConsumerServiceURL="https://sp.example.com/" Destination="https://localhost:8543/realms/demo/protocol/saml" ForceAuthn="false" ID="1cc6ba19-b34d-48d5-ac44-3c4b2598b6c3" IssueInstant="2025-12-10T16:29:41.177Z" ProtocolBinding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Version="2.0">
  <saml:Issuer>https://sp.example.com/</saml:Issuer>
  <samlp:NameIDPolicy Format="urn:oasis:names:tc:SAML:2.0:nameid-format:transient"/>
  <samlp:RequestedAuthnContext>
    <saml:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:TimeSyncToken</saml:AuthnContextClassRef>
    <saml:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport</saml:AuthnContextClassRef>
  </samlp:RequestedAuthnContext>
</samlp:AuthnRequest>
```

The request asks for a set of authentication contexts that will be evaluated in order by Keycloak until one of them can be satisfied. This means that Keycloak knows that URI declaration and there is a level defined in the mapping that is sufficient for the request. If none of the specified contexts can be satisfied, then Keycloak returns a `<Response>` message with a second-level `<StatusCode>` of `urn:oasis:names:tc:SAML:2.0:status:NoAuthnContext`.

```
<samlp:Response xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion" Destination="https://sp.example.com/" ID="ID_70625c64-004a-44cc-b338-74802f0cc37d" InResponseTo="b935f66a-8185-40b6-990c-bd5f4971773e" IssueInstant="2025-12-10T17:00:46.894Z" Version="2.0">
  <saml:Issuer>https://localhost:8543/auth/realms/demo</saml:Issuer>
  <samlp:Status>
    <samlp:StatusCode Value="urn:oasis:names:tc:SAML:2.0:status:Responder">
      <samlp:StatusCode Value="urn:oasis:names:tc:SAML:2.0:status:NoAuthnContext"/>
    </samlp:StatusCode>
  </samlp:Status>
</samlp:Response>
```

If Keycloak can satisfy the request, the step-up authentication flow is called to authenticate the user with the mapped level. The SAML response, using the `<AuthnStatement>` element inside the successful response, will contain again the `<AuthnContextClassRef>` with the context finally applied.

```
<samlp:Response xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion" Destination="https://sp.example.com/" ID="ID_079be770-ace4-4a0b-bdad-2042267a9bb9" InResponseTo="1cc6ba19-b34d-48d5-ac44-3c4b2598b6c3" IssueInstant="2025-12-10T16:29:41.228Z" Version="2.0">
  <saml:Issuer>https://localhost:8543/auth/realms/demo</saml:Issuer>
  <samlp:Status>
    <samlp:StatusCode Value="urn:oasis:names:tc:SAML:2.0:status:Success"/>
  </samlp:Status>
  <saml:Assertion xmlns="urn:oasis:names:tc:SAML:2.0:assertion" ID="ID_0278d50c-03c2-4ae0-accd-16161577faa5" IssueInstant="2025-12-10T16:29:41.228Z" Version="2.0">
    <saml:Issuer>https://localhost:8543/auth/realms/demo</saml:Issuer>
    <saml:Subject>
      <saml:NameID Format="urn:oasis:names:tc:SAML:2.0:nameid-format:transient">G-b1e11e8a-c002-47bc-8060-64f6f11fe267</saml:NameID>
      <saml:SubjectConfirmation Method="urn:oasis:names:tc:SAML:2.0:cm:bearer">
        <saml:SubjectConfirmationData InResponseTo="1cc6ba19-b34d-48d5-ac44-3c4b2598b6c3" NotOnOrAfter="2025-12-10T16:34:39.228Z" Recipient="https://sp.example.com/"/>
      </saml:SubjectConfirmation>
    </saml:Subject>
    <saml:Conditions NotBefore="2025-12-10T16:29:39.228Z" NotOnOrAfter="2025-12-10T16:30:39.228Z">
      <saml:AudienceRestriction>
        <saml:Audience>https://sp.example.com/</saml:Audience>
      </saml:AudienceRestriction>
    </saml:Conditions>
    <saml:AuthnStatement AuthnInstant="2025-12-10T16:29:41.229Z" SessionIndex="pZJ1Scsa1MyQ2KkesFEbedPq::56687d60-82d9-49bd-8609-dbc38153e55f" SessionNotOnOrAfter="2025-12-11T02:29:41.229Z">
      <saml:AuthnContext>
        <saml:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:TimeSyncToken</saml:AuthnContextClassRef>
      </saml:AuthnContext>
    </saml:AuthnStatement>
  </saml:Assertion>
</samlp:Response>
```

For more information about the processing rules defined in the SAML standard for the `<RequestedAuthnContext>` element, see Section 3.3.2.2.1 Element `<RequestedAuthnContext>` of the **Assertions and Protocols for the OASIS Security Assertion Markup Language (SAML) V2.0** (`saml-core-2.0-os.pdf` document). The complete SAML v2.0 OASIS Standard set (PDF format) and schema files are available in the [Security Assertion Markup Language (SAML) v2.0](https://www.oasis-open.org/standard/saml/) page. Note that the spec also defines a comparison (exact, minimum, better and maximum), that complicates what Level of Authentication will be finally selected in Keycloak.

When the `step-up-authentication-saml` feature is enabled, the [ACR to Level of Authentication (LoA) Mapping](#_mapping-acr-to-loa-realm) is a table with three values: the OpenID Connect ACR, the SAML URI context and the Level of Authentication (LoA). The mapping can be defined for both client types. It can also be overridden at client level, but, in this case, only the URI and the LoA are present. The minimum ACR value which is allowed for the client can also be defined in the configuration. The step-up authentication options for SAML clients are placed in the **Advanced** tab, section **Advanced Settings**. See chapter [Creating a SAML client](#_client-saml-configuration) for details.

In summary, when the step-up authentication is configured for SAML, Keycloak will process the specified context level in the SAML request and the minimum ACR allowed for the client (if they are present) to obtain the LoA (integer level) that should be reached in the authentication. If there is no available level that satisfies the request, the error is returned per specification. If the URI/LoA mapping returns a level that satisfies the request, the authentication flow is started, enforcing that Level of Authentication to be reached.

In order to maintain backwards compatibility, Keycloak does not return an error and continues adding the previous ACR `urn:oasis:names:tc:SAML:2.0:ac:classes:unspecified` to the response in the following situations:

- The new feature `step-up-authentication-saml` is not enabled in Keycloak.
- The SAML client does not define any mapping between context URI and LoA. Use the client mapping instead of the general realm mapping when you just need to apply the step-up for some specific SAML clients.
- The `AuthnContextClassRef mapper` is not executed. This mapper is provided by a new default client scope `AuthnContextClassRef` which is now added to SAML clients by default. It is in charge of adding the resulting `<AuthnContextClassRef>` to the response.
  
  For new realms created with the `step-up-authentication-saml` feature enabled, the mapper and the client scope `AuthnContextClassRef` are automatically created and assigned to SAML clients. But, for exiting realms, if you want to use this preview feature, the client scope and the mapper should be created and assigned to the client manually. When the feature is promoted to supported, the migration will also create the client scope for existing realms if the feature is not disabled at startup.
- The LoA calculated at request time is not achieved by the authentication flow. For example, if the authentication flow used for authentication is not a step-up flow, or there is a misconfiguration between the URI/LoA mapping and the final level reached in the step-up authentication flow.

#### [](#_registration-rc-client-flows)Registration or Reset credentials requested by client

Usually when the user is redirected to the Keycloak from client application, the `browser` flow is triggered. This flow may allow the user to [register](#con-user-registration_server_administration_guide) in case that realm registration is enabled and the user clicks `Register` on the login screen. Also, if [Forget password](#enabling-forgot-password) is enabled for the realm, the user can click `Forget password` on the login screen, which triggers the `Reset credentials` flow where users can reset credentials after email address confirmation.

Sometimes it can be useful for the client application to directly redirect the user to the **Registration** screen or to the **Reset credentials** flow. The resulting action will match the action of when the user clicks **Register** or **Forget password** on the normal login screen. Automatic redirect to the registration or reset-credentials screen can be done as follows:

- When the client wants the user to be redirected directly to the registration, the OIDC client should add parameter `prompt=create` to the login request. As a deprecated alternative, clients can replace the very last snippet from the OIDC login URL path (`/auth`) with `/registrations` . So the full URL might be similar to the following: `https://keycloak.example.com/realms/your_realm/protocol/openid-connect/registrations`. The `prompt=create` is recommended as it is [a specification standard](https://openid.net/specs/openid-connect-prompt-create-1_0.html).
- When the client wants a user to be redirected directly to the `Reset credentials` flow, the OIDC client should replace the very last snippet from the OIDC login URL path (`/auth`) with `/forgot-credentials` .

The preceding steps are the only supported method for a client to directly request a registration or reset-credentials flow. For security purposes, it is not supported and recommended for client applications to bypass OIDC/SAML flows and directly redirect to other Keycloak endpoints (such as for instance endpoints under `/realms/realm_name/login-actions` or `/realms/realm_name/broker`).

### [](#_user_session_limits)User session limits

Limits on the number of session that a user can have can be configured. Sessions can be limited per realm or per client.

To add session limits to a flow, perform the following steps.

01. Click **Add step** for the flow.
02. Select **User session count limiter** from the item list.
03. Click **Add**.
04. Click **Required** for the **User Session Count Limiter** authentication type to set its requirement to required.
05. Click **⚙️** (gear icon) for the **User Session Count Limiter**.
06. Enter an alias for this config.
07. Enter the required maximum number of sessions that a user can have in this realm. For example, if 2 is the value, 2 SSO sessions is the maximum that each user can have in this realm. If 0 is the value, this check is disabled.
08. Enter the required maximum number of sessions a user can have for the client. For example, if 2 is the value, then 2 SSO sessions is the maximum in this realm for each client. So when a user is trying to authenticate to client `foo`, but that user has already authenticated in 2 SSO sessions to client `foo`, either the authentication will be denied or an existing sessions will be killed based on the behavior configured. If a value of 0 is used, this check is disabled. If both session limits and client session limits are enabled, it makes sense to have client session limits to be always lower than session limits. The limit per client can never exceed the limit of all SSO sessions of this user.
09. Select the behavior that is required when the user tries to create a session after the limit is reached. Available behaviors are:
    
    - **Deny new session** - when a new session is requested and the session limit is reached, no new sessions can be created.
    - **Terminate oldest session** - when a new session is requested and the session limit has been reached, the oldest session will be removed and the new session created.
10. Optionally, add a custom error message to be displayed when the limit is reached.

Note that the user session limits should be added to your bound **Browser flow**, **Direct grant flow**, **Reset credentials** and also to any **Post broker login flow**. The authenticator should be added at the point when the user is already known during authentication (usually at the end of the authentication flow) and should be typically REQUIRED. Note that it is not possible to have ALTERNATIVE and REQUIRED executions at the same level.

For most of authenticators like `Direct grant flow`, `Reset credentials` or `Post broker login flow`, it is recommended to add the authenticator as REQUIRED at the end of the authentication flow. Here is an example for the `Reset credentials` flow:

![Authentication User Session Limits Reset Credentials Flow](./images/authentication-user-session-limits-resetcred.png)

For `Browser` flow, consider not adding the Session Limits authenticator at the top level flow. This recommendation is due to the `Cookie` authenticator, which automatically re-authenticates users based on SSO cookie. It is at the top level and it is better to not check session limits during SSO re-authentication because a user session already exists. So instead, consider adding a separate ALTERNATIVE subflow, such as the following `authenticate-user-with-session-limit` example at the same level like `Cookie`. Then you can add a REQUIRED subflow, in the following ``real-authentication-subflow`example, as a nested subflow of `authenticate-user-with-session-limit`` and add a `User Session Limit` at the same level as well. Inside the `real-authentication-subflow`, you can add real authenticators in a similar fashion to the default browser flow. The following example flow allows to users to authenticate with an identity provider or with password and OTP:

![Authentication User Session Limits Browser Flow](./images/authentication-user-session-limits-browser.png)

Regarding `Post Broker login flow`, you can add the `User Session Limits` as the only authenticator in the authentication flow as long as you have no other authenticators that you trigger after authentication with your identity provider. However, make sure that this flow is configured as `Post Broker Flow` at your identity providers. This requirement exists needed so that the authentication with Identity providers also participates in the session limits.

Currently, the administrator is responsible for maintaining consistency between the different configurations. So make sure that all your flows use same the configuration of `User Session Limits`.

User session limit feature is not available for CIBA.

### [](#script-authenticator)Script Authenticator

Ability to upload scripts through the Admin Console and REST endpoints is deprecated.

For more details see [JavaScript Providers](https://www.keycloak.org/docs/26.6.3/server_development/#_script_providers).

### [](#_kerberos)Kerberos

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/kerberos.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fkerberos.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fkerberos.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak supports login with a Kerberos ticket through the Simple and Protected GSSAPI Negotiation Mechanism (SPNEGO) protocol. SPNEGO authenticates transparently through the web browser after the user authenticates the session. For non-web cases, or when a ticket is not available during login, Keycloak supports login with Kerberos username and password.

A typical use case for web authentication is the following:

1. The user logs into the desktop.
2. The user accesses a web application secured by Keycloak using a browser.
3. The application redirects to Keycloak login.
4. Keycloak renders the HTML login screen with status 401 and HTTP header `WWW-Authenticate: Negotiate`
5. If the browser has a Kerberos ticket from desktop login, the browser transfers the desktop sign-on information to Keycloak in header `Authorization: Negotiate 'spnego-token'`. Otherwise, it displays the standard login screen, and the user enters the login credentials.
6. Keycloak validates the token from the browser and authenticates the user.
7. If using LDAPFederationProvider with Kerberos authentication support, Keycloak provisions user data from LDAP. If using KerberosFederationProvider, Keycloak lets the user update the profile and pre-fill login data.
8. Keycloak returns to the application. Keycloak and the application communicate through OpenID Connect or SAML messages. Keycloak acts as a broker to Kerberos/SPNEGO login. Therefore Keycloak authenticating through Kerberos is hidden from the application.

The [Negotiate](https://www.ietf.org/rfc/rfc4559.txt) www-authenticate scheme allows NTLM as a fallback to Kerberos and on some web browsers in Windows NTLM is supported by default. If a www-authenticate challenge comes from a server outside a browsers permitted list, users may encounter an NTLM dialog prompt. A user would need to click the cancel button on the dialog to continue as Keycloak does not support this mechanism. This situation can happen if Intranet web browsers are not strictly configured or if Keycloak serves users in both the Intranet and Internet. A [custom authenticator](https://github.com/keycloak/keycloak/issues/8989) can be used to restrict Negotiate challenges to a whitelist of hosts.

Perform the following steps to set up Kerberos authentication:

1. The setup and configuration of the Kerberos server (KDC).
2. The setup and configuration of the Keycloak server.
3. The setup and configuration of the client machines.

#### [](#setup-of-kerberos-server)Setup of Kerberos server

The steps to set up a Kerberos server depends on the operating system (OS) and the Kerberos vendor. Consult Windows Active Directory, MIT Kerberos, and your OS documentation for instructions on setting up and configuring a Kerberos server.

During setup, perform these steps:

1. Add some user principals to your Kerberos database. You can also integrate your Kerberos with LDAP, so user accounts provision from the LDAP server.
2. Add service principal for "HTTP" service. For example, if the Keycloak server runs on `www.mydomain.org`, add the service principal `HTTP/www.mydomain.org@<kerberos realm>`.
   
   On MIT Kerberos, you run a "kadmin" session. On a machine with MIT Kerberos, you can use the command:

```
sudo kadmin.local
```

Then, add HTTP principal and export its key to a keytab file with commands such as:

```
addprinc -randkey HTTP/www.mydomain.org@MYDOMAIN.ORG
ktadd -k /tmp/http.keytab HTTP/www.mydomain.org@MYDOMAIN.ORG
```

Ensure the keytab file `/tmp/http.keytab` is accessible on the host where Keycloak is running.

#### [](#setup-and-configuration-of-keycloak-server)Setup and configuration of Keycloak server

Install a Kerberos client on your machine.

Procedure

1. Install a Kerberos client. If your machine runs Fedora, Ubuntu, or RHEL, install the [freeipa-client](https://www.freeipa.org/page/Downloads) package, containing a Kerberos client and other utilities.
2. Configure the Kerberos client (on Linux, the configuration settings are in the [/etc/krb5.conf](https://web.mit.edu/kerberos/krb5-1.21/doc/admin/conf_files/krb5_conf.html) file ).
   
   Add your Kerberos realm to the configuration and configure the HTTP domains your server runs on.
   
   For example, for the MYDOMAIN.ORG realm, you can configure the `domain_realm` section like this:
   
   ```
   [domain_realm]
     .mydomain.org = MYDOMAIN.ORG
     mydomain.org = MYDOMAIN.ORG
   ```
3. Export the keytab file with the HTTP principal and ensure the file is accessible to the process running the Keycloak server. For production, ensure that the file is readable by this process only.
   
   For the MIT Kerberos example above, we exported keytab to the `/tmp/http.keytab` file. If your *Key Distribution Centre (KDC)* and Keycloak run on the same host, the file is already available.

##### [](#enabling-spnego-processing)Enabling SPNEGO processing

By default, Keycloak disables SPNEGO protocol support. To enable it, go to the [browser flow](#_authentication-flows) and enable **Kerberos**.

Browser flow

![Browser Flow](./images/browser-flow.png)

Set the **Kerberos** requirement from *disabled* to *alternative* (Kerberos is optional) or *required* (browser must have Kerberos enabled). If you have not configured the browser to work with SPNEGO or Kerberos, Keycloak falls back to the regular login screen.

##### [](#configure-kerberos-user-storage-federation-providers)Configure Kerberos user storage federation providers

You must now use [User Storage Federation](#_user-storage-federation) to configure how Keycloak interprets Kerberos tickets. Two different federation providers exist with Kerberos authentication support.

To authenticate with Kerberos backed by an LDAP server, configure the [LDAP Federation Provider](#_ldap).

Procedure

1. Go to the configuration page for your LDAP provider.
   
   Ldap kerberos integration
   
   ![LDAP Kerberos Integration](./images/ldap-kerberos.png)
2. Toggle **Allow Kerberos authentication** to **ON**

**Allow Kerberos authentication** makes Keycloak use the Kerberos principal access user information so information can import into the Keycloak environment.

If an LDAP server is not backing up your Kerberos solution, use the **Kerberos** User Storage Federation Provider.

Procedure

1. Click **User Federation** in the menu.
2. Select **Kerberos** from the **Add provider** select box.
   
   Kerberos user storage provider
   
   ![Kerberos User Storage Provider](./images/kerberos-provider.png)

The **Kerberos** provider parses the Kerberos ticket for simple principal information and imports the information into the local Keycloak database. User profile information, such as first name, last name, and email, are not provisioned.

#### [](#setup-and-configuration-of-client-machines)Setup and configuration of client machines

Client machines must have a Kerberos client and set up the `krb5.conf` as described [above](#_server_setup). The client machines must also enable SPNEGO login support in their browser. See [configuring Firefox for Kerberos](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/7/html/system-level_authentication_guide/configuring_applications_for_sso) if you are using the Firefox browser.

The `.mydomain.org` URI must be in the `network.negotiate-auth.trusted-uris` configuration option.

In Windows domains, clients do not need to adjust their configuration. Internet Explorer and Edge can already participate in SPNEGO authentication.

#### [](#example-setups)Example setups

##### [](#keycloak-and-freeipa-docker-image)Keycloak and FreeIPA docker image

When you install [docker](https://www.docker.com/), run a docker image with the FreeIPA server installed. FreeIPA provides an integrated security solution with MIT Kerberos and 389 LDAP server. The image also contains a Keycloak server configured with an LDAP Federation provider and enabled SPNEGO/Kerberos authentication against the FreeIPA server. See details [here](https://github.com/mposolda/keycloak-freeipa-docker/blob/master/README.md).

##### [](#apacheds-testing-kerberos-server)ApacheDS testing Kerberos server

For quick testing and unit tests, use a simple [ApacheDS](https://directory.apache.org/apacheds/) Kerberos server. You must build Keycloak from the source and then run the Kerberos server with the maven-exec-plugin from our test suite. See details [here](https://github.com/keycloak/keycloak/blob/main/docs/tests.md#kerberos-server).

#### [](#credential-delegation)Credential delegation

Kerberos supports the credential delegation. Applications may need access to the Kerberos ticket so they can reuse it to interact with other services secured by Kerberos. Because the Keycloak server processed the SPNEGO protocol, you must propagate the GSS credential to your application within the OpenID Connect token claim or a SAML assertion attribute. Keycloak transmits this to your application from the Keycloak server. To insert this claim into the token or assertion, each application must enable the built-in protocol mapper `gss delegation credential`. This mapper is available in the **Mappers** tab of the application’s client page. See [Protocol Mappers](#_protocol-mappers) chapter for more details.

Applications must deserialize the claim it receives from Keycloak before using it to make GSS calls against other services. When you deserialize the credential from the access token to the GSSCredential object, create the GSSContext with this credential passed to the `GSSManager.createContext` method. For example:

```
// Obtain accessToken in your application.
KeycloakPrincipal keycloakPrincipal = (KeycloakPrincipal) servletReq.getUserPrincipal();
AccessToken accessToken = keycloakPrincipal.getKeycloakSecurityContext().getToken();

// Retrieve Kerberos credential from accessToken and deserialize it
String serializedGssCredential = (String) accessToken.getOtherClaims().
    get(org.keycloak.common.constants.KerberosConstants.GSS_DELEGATION_CREDENTIAL);

GSSCredential deserializedGssCredential = org.keycloak.common.util.KerberosSerializationUtils.
    deserializeCredential(serializedGssCredential);

// Create GSSContext to call other Kerberos-secured services
GSSContext context = gssManager.createContext(serviceName, krb5Oid,
    deserializedGssCredential, GSSContext.DEFAULT_LIFETIME);
```

Configure `forwardable` Kerberos tickets in `krb5.conf` file and add support for delegated credentials to your browser.

Credential delegation has security implications, so use it only if necessary and only with HTTPS. See [this article](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/7/html/system-level_authentication_guide/configuring_applications_for_sso) for more details and an example.

#### [](#cross-realm-trust)Cross-realm trust

In the Kerberos protocol, the `realm` is a set of Kerberos principals. The definition of these principals exists in the Kerberos database, which is typically an LDAP server.

The Kerberos protocol allows cross-realm trust. For example, if 2 Kerberos realms, A and B, exist, then cross-realm trust will allow the users from realm A to access realm B’s resources. Realm B trusts realm A.

Kerberos cross-realm trust

![kerberos trust basic](./images/kerberos-trust-basic.png)

The Keycloak server supports cross-realm trust. To implement this, perform the following:

- Configure the Kerberos servers for the cross-realm trust. Implementing this step depends on the Kerberos server implementations. This step is necessary to add the Kerberos principal `krbtgt/B@A` to the Kerberos databases of realm A and B. This principal must have the same keys on both Kerberos realms. The principals must have the same password, key version numbers, and ciphers in both realms. Consult the Kerberos server documentation for more details.

The cross-realm trust is unidirectional by default. You must add the principal `krbtgt/A@B` to both Kerberos databases for bidirectional trust between realm A and realm B. However, trust is transitive by default. If realm B trusts realm A and realm C trusts realm B, then realm C trusts realm A without the principal, `krbtgt/C@A`, available. Additional configuration (for example, `capaths`) may be necessary on the Kerberos client-side so clients can find the trust path. Consult the Kerberos documentation for more details.

- Configure Keycloak server
  
  - When using an LDAP storage provider with Kerberos support, configure the server principal for realm B, as in this example: `HTTP/mydomain.com@B`. The LDAP server must find the users from realm A if users from realm A are to successfully authenticate to Keycloak, because Keycloak must perform the SPNEGO flow and then find the users.

Finding users is based on the LDAP storage provider option `Kerberos principal attribute`. When this is configured for instance with value like `userPrincipalName`, then after SPNEGO authentication of user `john@A`, Keycloak will try to lookup LDAP user with attribute `userPrincipalName` equivalent to `john@A`. If `Kerberos principal attribute` is left empty, then Keycloak will lookup the LDAP user based on the prefix of his kerberos principal with the realm omitted. For example, Kerberos principal user `john@A` must be available in the LDAP under username `john`, so typically under an LDAP DN such as `uid=john,ou=People,dc=example,dc=com`. If you want users from realm A and B to authenticate, ensure that LDAP can find users from both realms A and B.

- When using a Kerberos user storage provider (typically, Kerberos without LDAP integration), configure the server principal as `HTTP/mydomain.com@B`, and users from Kerberos realms A and B must be able to authenticate.

Users from multiple Kerberos realms are allowed to authenticate as every user would have attribute `KERBEROS_PRINCIPAL` referring to the kerberos principal used for authentication and this is used for further lookups of this user. To avoid conflicts when there is user `john` in both kerberos realms `A` and `B`, the username of the Keycloak user might contain the kerberos realm lowercased. For instance username would be `john@a`. Just in case when realm matches with the configured `Kerberos realm`, the realm suffix might be omitted from the generated username. For instance username would be `john` for the Kerberos principal `john@A` as long as the `Kerberos realm` is configured on the Kerberos provider is `A`.

#### [](#troubleshooting)Troubleshooting

If you have issues, enable additional logging to debug the problem:

- Enable `Debug` flag in the Admin Console for Kerberos or LDAP federation providers
- Enable TRACE logging for category `org.keycloak` to receive more information in server logs
- Add system properties `-Dsun.security.krb5.debug=true` and `-Dsun.security.spnego.debug=true`

### [](#_x509)X.509 client certificate user authentication

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/x509.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fx509.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fx509.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak supports logging in with an X.509 client certificate if you have configured the server to use mutual SSL authentication.

A typical workflow:

- A client sends an authentication request over SSL/TLS channel.
- During the SSL/TLS handshake, the server and the client exchange their x.509/v3 certificates.
- The container (WildFly) validates the certificate PKIX path and the certificate expiration date.
- The x.509 client certificate authenticator validates the client certificate by using the following methods:
  
  - Checks the certificate revocation status by using CRL or CRL Distribution Points.
  - Checks the Certificate revocation status by using OCSP (Online Certificate Status Protocol).
  - Validates whether the key in the certificate matches the expected key.
  - Validates whether the extended key in the certificate matches the expected extended key.
- If any of the these checks fail, the x.509 authentication fails. Otherwise, the authenticator extracts the certificate identity and maps it to an existing user.

When the certificate maps to an existing user, the behavior diverges depending on the authentication flow:

- In the Browser Flow, the server prompts users to confirm their identity or sign in with a username and password.
- In the Direct Grant Flow, the server signs in the user.

Note that it is the responsibility of the web container to validate certificate PKIX path. X.509 authenticator on the Keycloak side provides just the additional support for check the certificate expiration, certificate revocation status and key usage. If you are using Keycloak deployed behind reverse proxy, make sure that your reverse proxy is configured to validate PKIX path. If you do not use reverse proxy and users directly access the WildFly, you should be fine as WildFly makes sure that PKIX path is validated as long as it is configured as described below.

#### [](#features-2)Features

Supported Certificate Identity Sources:

- Match SubjectDN by using regular expressions
- X500 Subject’s email attribute
- X500 Subject’s email from Subject Alternative Name Extension (RFC822Name General Name)
- X500 Subject’s other name from Subject Alternative Name Extension. This other name is the User Principal Name (UPN), typically.
- X500 Subject’s Common Name attribute
- Match IssuerDN by using regular expressions
- Certificate Serial Number
- Certificate Serial Number and IssuerDN
- SHA-256 Certificate thumbprint
- Full certificate in PEM format

##### [](#regular-expressions)Regular expressions

Keycloak extracts the certificate identity from Subject DN or Issuer DN by using a regular expression as a filter. For example, this regular expression matches the email attribute:

```
emailAddress=(.*?)(?:,|$)
```

The regular expression filtering applies if the `Identity Source` is set to either `Match SubjectDN using regular expression` or `Match IssuerDN using regular expression`.

###### [](#mapping-certificate-identity-to-an-existing-user)Mapping certificate identity to an existing user

The certificate identity mapping can map the extracted user identity to an existing user’s username, email, or a custom attribute whose value matches the certificate identity. For example, setting `Identity source` to *Subject’s email* or `User mapping method` to *Username or email* makes the X.509 client certificate authenticator use the email attribute in the certificate’s Subject DN as the search criteria when searching for an existing user by username or by email.

- If you disable **Login with email** at realm settings, the same rules apply to certificate authentication. Users are unable to log in by using the email attribute.
- Using `Certificate Serial Number and IssuerDN` as an identity source requires two custom attributes for the serial number and the IssuerDN.
- `SHA-256 Certificate thumbprint` is the lowercase hexadecimal representation of SHA-256 certificate thumbprint.
- Using `Full certificate in PEM format` as an identity source is limited to the custom attributes mapped to external federation sources, such as LDAP. Keycloak cannot store certificates in its database due to length limitations, so in the case of LDAP, you must enable `Always Read Value From LDAP`.

###### [](#extended-certificate-validation)Extended certificate validation

- Revocation status checking using CRL.
- Revocation status checking using CRL/Distribution Point.
- Revocation status checking using OCSP/Responder URI.
- Certificate KeyUsage validation.
- Certificate ExtendedKeyUsage validation.

#### [](#_browser_flow)Adding X.509 client certificate authentication to browser flows

01. Click **Authentication** in the menu.
02. Click the **Browser** flow.
03. From the **Action** list, select **Duplicate**.
04. Enter a name for the copy.
05. Click **Duplicate**.
06. Click **Add step**.
07. Click "X509/Validate Username Form".
08. Click **Add**.
    
    X509 execution
    
    ![X509 Execution](./images/x509-execution.png)
09. Click and drag the "X509/Validate Username Form" over the "Browser Forms" execution.
10. Set the requirement to "ALTERNATIVE".
    
    X509 browser flow
    
    ![X509 Browser Flow](./images/x509-browser-flow.png)
11. Click the **Action** menu.
12. Click the **Bind flow**.
13. Click the **Browser flow** from the drop-down list.
14. Click **Save**.
    
    X509 browser flow bindings
    
    ![X509 Browser Flow Bindings](./images/x509-browser-flow-bindings.png)

#### [](#_x509-config)Configuring X.509 client certificate authentication

X509 configuration

![X509 Configuration](./images/x509-configuration.png)

**User Identity Source**

Defines the method for extracting the user identity from a client certificate.

**Canonical DN representation enabled**

Defines whether to use the canonical format to determine a distinguished name. The format adds additional normalization rules to the RFC 2253 conformant string representation. The official [Java API documentation](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/javax/security/auth/x500/X500Principal.html#getName%28java.lang.String%29) describes the additions in detail. This option affects the following User Identity Sources only:

- *Match SubjectDN using regular expression*
- *Match IssuerDN using regular expression*
- *Certificate Serial Number and IssuerDN* (the IssuerDN portion)

Do not enable this option to retain backward compatibility with existing Keycloak instances.

The `canonical` format performs modifications over the presented DN in the certificate (e.g., leading and trailing spaces are removed or the entire name is converted to uppercase and lowercase using the English locale). This behavior can lead to collisions in user matching if the Certificate Authorities (CA) that issue the certificates do not validate the assigned DN (e.g., if DNs are issued using any locale and present problems when performing the upper and lower case in English). Do not enable this option if you cannot guarantee that the canonical representation avoids duplications for later matching.

**Enable Serial Number hexadecimal representation**

Represent the [serial number](https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.2) as hexadecimal. The serial number with the sign bit set to 1 must be left padded with 00 octet. For example, a serial number with decimal value *161*, or *a1* in hexadecimal representation is encoded as *00a1*, according to RFC 5280. See [RFC 5280, appendix-B](https://datatracker.ietf.org/doc/html/rfc5280#appendix-B) for more details.

**A regular expression**

A regular expression to use as a filter for extracting the certificate identity. The expression must contain a single group.

**User Mapping Method**

Defines the method to match the certificate identity with an existing user. *Username or email* searches for existing users by username or email. *Custom Attribute Mapper* searches for existing users with a custom attribute that matches the certificate identity. The name of the custom attribute is configurable.

**A name of user attribute**

A custom attribute whose value matches against the certificate identity. Use multiple custom attributes when attribute mapping is related to multiple values, For example, 'Certificate Serial Number and IssuerDN'.

**CRL Checking Enabled**

Check the revocation status of the certificate by using the Certificate Revocation List. The location of the list is defined in the **CRL file path** attribute.

**Enable CRL Distribution Point to check certificate revocation status**

Use CDP to check the certificate revocation status. Most PKI authorities include CDP in their certificates.

**CRL file path**

The path to a file containing a CRL list. The value must be a path to a valid file if the **CRL Checking Enabled** option is enabled.

**CRL abort if non updated**

A CRL conforming to [RFC 5280](https://datatracker.ietf.org/doc/html/rfc5280#section-5.1.2.5) contains a next update field that indicates the date by which the next CRL will be issued. When that time is passed, the CRL is considered outdated and it should be refreshed. If this option is `true`, the authentication will fail if the CRL is outdated (recommended). If the option is set to `false`, the outdated CRL is still used to validate the user certificates.

**OCSP Checking Enabled**

Checks the certificate revocation status by using Online Certificate Status Protocol.

**OCSP Fail-Open Behavior**

By default the OCSP check must return a positive response in order to continue with a successful authentication. Sometimes however this check can be inconclusive: for example, the OCSP server could be unreachable, overloaded, or the client certificate may not contain an OCSP responder URI. When this setting is turned ON, authentication will be denied only if an explicit negative response is received by the OCSP responder and the certificate is definitely revoked. If a valid OCSP response is not available the authentication attempt will be accepted.

OCSP retry behavior is configured server-wide through the HTTP client provider. See &lt;@links.server id="outgoinghttp"/&gt; for details on configuring retry settings for all outgoing HTTP requests, including OCSP validation.

**OCSP Responder URI**

Override the value of the OCSP responder URI in the certificate.

**Validate Key Usage**

Verifies the certificate’s KeyUsage extension bits are set. For example, "digitalSignature,KeyEncipherment" verifies if bits 0 and 2 in the KeyUsage extension are set. Leave this parameter empty to disable the Key Usage validation. See [RFC 5280, Section-4.2.1.3](https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.3) for more information. Keycloak raises an error when a key usage mismatch occurs.

**Validate Extended Key Usage**

Verifies one or more purposes defined in the Extended Key Usage extension. See [RFC 5280, Section-4.2.1.12](https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.12) for more information. Leave this parameter empty to disable the Extended Key Usage validation. Keycloak raises an error when flagged as critical by the issuing CA and a key usage extension mismatch occurs.

**Validate Certificate Policy**

Verifies one or more policy OIDs as defined in the Certificate Policy extension. See [RFC 5280, Section-4.2.1.4](https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.4). Leave the parameter empty to disable the Certificate Policy validation. Multiple policies should be separated using a comma.

**Certificate Policy Validation Mode**

When more than one policy is specified in the `Validate Certificate Policy` setting, it decides whether the matching should check for all requested policies to be present, or one match is enough for a successful authentication. Default value is `All`, meaning that all requested policies should be present in the client certificate.

**Bypass identity confirmation**

If enabled, X.509 client certificate authentication does not prompt the user to confirm the certificate identity. Keycloak signs in the user upon successful authentication.

**Revalidate client certificate**

If set, the client certificate trust chain will be always verified at the application level using the certificates present in the configured trust store. This can be useful if the underlying web server does not enforce client certificate chain validation, for example because it is behind a non-validating load balancer or reverse proxy, or when the number of allowed CAs is too large for the mutual SSL negotiation (most browsers cap the maximum SSL negotiation packet size at 32767 bytes, which corresponds to about 200 advertised CAs). By default this option is off.

#### [](#adding-x-509-client-certificate-authentication-to-a-direct-grant-flow)Adding X.509 Client Certificate Authentication to a Direct Grant Flow

01. Click **Authentication** in the menu.
02. Select **Duplicate** from the "Action list" to make a copy of the built-in "Direct grant" flow.
03. Enter a name for the copy.
04. Click **Duplicate**.
05. Click the created flow.
06. Click the trash can icon 🗑️ of the "Username Validation" and click **Delete**.
07. Click the trash can icon 🗑️ of the "Password" and click **Delete**.
08. Click **Add step**.
09. Click "X509/Validate Username".
10. Click **Add**.
    
    X509 direct grant execution
    
    ![X509 Direct Grant Execution](./images/x509-directgrant-execution.png)
11. Set up the x509 authentication configuration by following the steps described in the [x509 Browser Flow](#_browser_flow) section.
12. Click the **Bindings** tab.
13. Click the **Direct Grant Flow** drop-down list.
14. Click the newly created "x509 Direct Grant" flow.
15. Click **Save**.
    
    X509 direct grant flow bindings
    
    ![X509 Direct Grant Flow Bindings](./images/x509-directgrant-flow-bindings.png)

##### [](#example-using-curl)Example using CURL

The following example shows how to obtain an access token for a user in the realm `test` with the direct grant flow. The example is using **OAuth2 Resource Owner Password Credentials Grant** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section and the confidential client `resource-owner`:

```
curl \
  -d "client_id=resource-owner" \
  -d "client_secret=40cc097b-2a57-4c17-b36a-8fdf3fc2d578" \
  -d "grant_type=password" \
  --cacert /tmp/truststore.pem \
  --cert /tmp/keystore.pem:kssecret \
  "https://localhost:8543/realms/test/protocol/openid-connect/token"
```

The file `/tmp/truststore.pem` points to the file with the truststore containing the certificate of the Keycloak server. The file `/tmp/keystore.pem` contains the private key and certificates corresponding to the Keycloak user, which would be successfully authenticated by this request. It is dependent on the configuration of the authenticator on how exactly is the content from the certificate mapped to the Keycloak user as described in [the configuration section](#_x509-config). The `kssecret` might be the password of this keystore file.

According to your environment, it might be needed to use more options to CURL commands like for instance:

- Option `--insecure` if you are using self-signed certificates
- Option `--capath` to include the whole directory containing the certificate authority path
- Options `--cert-type` or `--key-type` in case you want to use different files than `PEM`

Please consult the documentation of the `curl` tool for the details if needed. If you are using other tools than `curl`, consult the documentation of your tool. However, the setup would be similar. A need exists to include keystore and truststore as well as client credentials in case you are using a confidential client.

If it is possible, it is preferred to use [Service accounts](#_service_accounts) together with the MTLS client authentication (client authenticator `X509 Certificate`) rather than using the Direct grant with X.509 authentication as direct grant may require sharing of the user certificate with client applications. When using service account, the tokens are obtained on behalf of the client itself, which in general is better and more secure practice.

### [](#webauthn_server_administration_guide)W3C Web Authentication (WebAuthn)

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/webauthn.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fwebauthn.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fwebauthn.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak provides support for [W3C Web Authentication (WebAuthn)](https://www.w3.org/TR/webauthn/). Keycloak works as a WebAuthn’s [Relying Party (RP)](https://www.w3.org/TR/webauthn/#webauthn-relying-party).

WebAuthn’s operations success depends on the user’s WebAuthn supporting authenticator, browser, and platform. Make sure your authenticator, browser, and platform support the WebAuthn specification.

WebAuthn’s specification uses a `user.id` to map a public key credential to a specific user account in the Relying Party. This user ID handle is an opaque byte sequence with a maximum size of 64 bytes. Keycloak passes the internal database ID to the registration, which in common users is an UUID of 36 characters. But, if the user is from a external user federation provider, the internal Keycloak ID is an [storage ID](https://www.keycloak.org/docs/26.6.3/server_development/#storage-ids) in the form `f:<provider-id>:<user-id>` that can exceed the 64 byte limitation. Please take this into account and use short IDs for the federation provider component and for the users coming from that provider when combining the Storage SPI and WebAuthn.

#### [](#setup)Setup

The setup procedure of WebAuthn support for 2FA is the following:

##### [](#_webauthn-register)Check WebAuthn authenticator registration is enabled

1. Click **Authentication** in the menu.
2. Click the **Required Actions** tab.
3. Check action **Webauthn Register** switch is set to **ON**.

Toggle the **Default Action** switch to **ON** if you want all new users to be required to register their WebAuthn credentials.

#### [](#_webauthn-authenticator-setup)Enable WebAuthn authentication in the default browser flow

1. Click **Authentication** in the menu.
2. Click the **Browser** flow.
3. Locate the execution **WebAuthn Authenticator** inside the **Browser - Conditional 2FA** sub-flow.
4. Change the *requirement* from *Disabled* to *Alternative* for that execution.
   
   WebAuthn browser flow conditional with OTP
   
   ![WebAuthn browser flow conditional with OTP](./images/webauthn-browser-flow-conditional-with-OTP.png)

With this configuration, the users can choose between using WebAuthn and OTP for the second factor. As the sub-flow is *conditional*, they are only asked to present a 2FA credential (OTP or WebAuthn) if they have already registered one of the respective credential types. If a user has configured both credential types, the credential with the highest priority will be displayed by default. However, the **Try Another Way** option will appear so that the user has the alternative methods to log in.

If you want to substitute OTP for WebAuthn and maintain it as conditional:

1. Change *requirement* in **OTP Form** to *Disabled*.
2. Change *requirement* in **WebAuthn Authenticator** to *Alternative*.
   
   Webauthn browser flow conditional
   
   ![Webauthn browser flow conditional](./images/webauthn-browser-flow-conditional.png)

If you require WebAuthn for all users and enforce them to configure the credential if not configured:

1. Change *requirement* in **Browser - Conditional 2FA** to *Required*.
2. Change *requirement* in **OTP Form** to *Disabled*.
3. Change *requirement* in **WebAuthn Authenticator** to *Required*.
   
   Webauthn browser flow required
   
   ![Webauthn browser flow required](./images/webauthn-browser-flow-required.png)

You can see more examples of 2FA configurations in [2FA conditional workflow examples](#twofa-conditional-workflow-examples).

#### [](#authenticate-with-webauthn-authenticator)Authenticate with WebAuthn authenticator

After registering a WebAuthn authenticator, the user carries out the following operations:

- Open the login form. The user must authenticate with a username and password.
- The user’s browser asks the user to authenticate by using their WebAuthn authenticator.

#### [](#managing-webauthn-as-an-administrator)Managing WebAuthn as an administrator

##### [](#managing-credentials)Managing credentials

Keycloak manages WebAuthn credentials similarly to other credentials from [User credential management](#ref-user-credentials_server_administration_guide):

- Keycloak assigns users a required action to create a WebAuthn credential from the **Reset Actions** list and select **Webauthn Register**.
- Administrators can delete a WebAuthn credential by clicking **Delete**.
- Administrators can view the credential’s data, such as the AAGUID, by selecting **Show data…​**.
- Administrators can set a label for the credential by setting a value in the **User Label** field and saving the data.

##### [](#_webauthn-policy)Managing policy

Administrators can configure WebAuthn related operations as **WebAuthn Policy** per realm.

Procedure

1. Click **Authentication** in the menu.
2. Click the **Policy** tab.
3. Click the **WebAuthn Policy** tab.
4. Configure the items within the policy (see description below).
5. Click **Save**.

The configurable items and their description are as follows:

  Configuration Description

Relying Party Entity Name

The readable server name as a WebAuthn Relying Party. This item is mandatory and applies to the registration of the WebAuthn authenticator. The default setting is "keycloak". For more details, see [WebAuthn Specification](https://www.w3.org/TR/webauthn/#dictionary-pkcredentialentity).

Signature Algorithms

The algorithms telling the WebAuthn authenticator which signature algorithms to use for the [Public Key Credential](https://www.w3.org/TR/webauthn/#iface-pkcredential). Keycloak uses the Public Key Credential to sign and verify [Authentication Assertions](https://www.w3.org/TR/webauthn/#authentication-assertion). If no algorithms exist, the default [ES256](https://datatracker.ietf.org/doc/html/rfc8152#section-8.1) and [RS256](https://datatracker.ietf.org/doc/html/rfc7518#section-3.1) is adapted. ES256 and RS256 are an optional configuration item applying to the registration of WebAuthn authenticators. For more details, see [WebAuthn Specification](https://www.w3.org/TR/webauthn/#dictdef-publickeycredentialparameters).

Relying Party ID

The ID of a WebAuthn Relying Party that determines the scope of [Public Key Credentials](https://www.w3.org/TR/webauthn/#public-key-credential). The ID must be the origin’s effective domain. This ID is an optional configuration item applied to the registration of WebAuthn authenticators. If this entry is blank, Keycloak adapts the host part of Keycloak’s base URL. For more details, see [WebAuthn Specification](https://www.w3.org/TR/webauthn/).

Attestation Conveyance Preference

The WebAuthn API implementation on the browser ([WebAuthn Client](https://www.w3.org/TR/webauthn/#webauthn-client)) is the preferential method to generate Attestation statements. This preference is an optional configuration item applying to the registration of the WebAuthn authenticator. If no option exists, its behavior is the same as selecting "none". For more details, see [WebAuthn Specification](https://www.w3.org/TR/webauthn/).

Authenticator Attachment

The acceptable attachment pattern of a WebAuthn authenticator for the WebAuthn Client. This pattern is an optional configuration item applying to the registration of the WebAuthn authenticator. For more details, see [WebAuthn Specification](https://www.w3.org/TR/webauthn/#enumdef-authenticatorattachment).

Require Discoverable Credential

The option requiring that the WebAuthn authenticator generates the Public Key Credential as [Client-side discoverable Credential](https://www.w3.org/TR/webauthn-3/). This option applies to the registration of the WebAuthn authenticator. If left blank, its behavior is the same as selecting "No". For more details, see [WebAuthn Specification](https://www.w3.org/TR/webauthn/#dom-authenticatorselectioncriteria-requireresidentkey).

User Verification Requirement

The option requiring that the WebAuthn authenticator confirms the verification of a user. This is an optional configuration item applying to the registration of a WebAuthn authenticator and the authentication of a user by a WebAuthn authenticator. If no option exists, its behavior is the same as selecting "preferred". For more details, see [WebAuthn Specification for registering a WebAuthn authenticator](https://www.w3.org/TR/webauthn/#dom-authenticatorselectioncriteria-userverification) and [WebAuthn Specification for authenticating the user by a WebAuthn authenticator](https://www.w3.org/TR/webauthn/#dom-publickeycredentialrequestoptions-userverification).

Timeout

The timeout value, in seconds, for registering a WebAuthn authenticator and authenticating the user by using a WebAuthn authenticator. If set to zero, its behavior depends on the WebAuthn authenticator’s implementation. The default value is 0. For more details, see [WebAuthn Specification for registering a WebAuthn authenticator](https://www.w3.org/TR/webauthn/#dom-publickeycredentialcreationoptions-timeout) and [WebAuthn Specification for authenticating the user by a WebAuthn authenticator](https://www.w3.org/TR/webauthn/#dom-publickeycredentialrequestoptions-timeout).

Avoid Same Authenticator Registration

If enabled, Keycloak cannot re-register an already registered WebAuthn authenticator.

Acceptable AAGUIDs

The list of allowed AAGUIDs which a WebAuthn authenticator must register against. An AAGUID (Authenticator Attestation Global Unique Identifier) is a 128-bit identifier indicating the authenticator’s type (e.g., make and model). This option needs the **Attestation conveyance preference** to be configured (normally `Direct`) to ensure a trusted AAGUID is passed. Default attestation `None` is not reliable, and can anonymize the AAGUID to zero value. If you setup **Acceptable AAGUIDs** only those authenticators are valid to register new WebAuthn credentials.

#### [](#attestation-statement-verification)Attestation statement verification

When registering a WebAuthn authenticator, Keycloak verifies the trustworthiness of the attestation statement generated by the WebAuthn authenticator. Keycloak requires the trust anchor’s certificates imported into the [truststore](https://www.keycloak.org/server/keycloak-truststore).

To omit this validation, disable this truststore or set the WebAuthn policy’s configuration item "Attestation Conveyance Preference" to "none".

#### [](#managing-webauthn-credentials-as-a-user)Managing WebAuthn credentials as a user

##### [](#register-webauthn-authenticator)Register WebAuthn authenticator

The appropriate method to register a WebAuthn authenticator depends on whether the user has already registered an account on Keycloak.

##### [](#new-user)New user

If the **WebAuthn Register** required action is **Default Action** in a realm, new users must set up the Passkey after their first login.

Procedure

1. Open the login form.
2. Click **Register**.
3. Fill in the items on the form.
4. Click **Register**.

After successfully registering, the browser asks the user to enter the text of their WebAuthn authenticator’s label.

##### [](#existing-user)Existing user

If `WebAuthn Authenticator` is set up as required as shown in the first example, then when existing users try to log in, they are required to register their WebAuthn authenticator automatically:

Procedure

1. Open the login form.
2. Enter the items on the form.
3. Click **Save**.
4. Click **Login**.

After successful registration, the user’s browser asks the user to enter the text of their WebAuthn authenticator’s label.

#### [](#_webauthn_aia)Registering WebAuthn credentials using AIA

WebAuthn credentials can also be registered for a user using [Application Initiated Actions (AIA)](#con-aia_server_administration_guide). The actions **Webauthn Register** (`kc_action=webauthn-register`) and **Webauthn Register Passwordless** (`kc_action=webauthn-register-passwordless`) are available for the applications if enabled in the [Required actions tab](#proc-setting-default-required-actions_server_administration_guide).

Both required actions allow a parameter **skip\_if\_exists** that allows to skip the AIA execution if the user already has a credential of that type. The `kc_action_status` will be **success** if skipped. For example, adding the option to the common WebAuthn register action is just using the following query parameter `kc_action=webauthn-register:skip_if_exists`.

#### [](#_webauthn_passwordless)Passwordless WebAuthn together with Two-Factor

Keycloak uses WebAuthn for two-factor authentication, but you can use WebAuthn as the first-factor authentication. In this case, users with `passwordless` WebAuthn credentials can authenticate to Keycloak without a password. Keycloak can use WebAuthn as both the passwordless and two-factor authentication mechanism in the context of a realm and a single authentication flow.

An administrator typically requires that Passkeys registered by users for the WebAuthn passwordless authentication meet different requirements. For example, the Passkeys may require users to authenticate to the Passkey using a PIN, or the Passkey attests with a stronger certificate authority.

Because of this, Keycloak permits administrators to configure a separate `WebAuthn Passwordless Policy`. There is a required `Webauthn Register Passwordless` action of type and separate authenticator of type `WebAuthn Passwordless Authenticator`.

##### [](#setup-2)Setup

Set up WebAuthn passwordless support as follows:

1. (if not already present) Register a new required action for WebAuthn passwordless support. Use the steps described in [Enable WebAuthn Authenticator Registration](#_webauthn-register). Register the `Webauthn Register Passwordless` action.
2. Configure the policy. You can use the steps and configuration options described in [Managing Policy](#_webauthn-policy). Perform the configuration in the Admin Console in the tab **WebAuthn Passwordless Policy**. Typically the requirements for the Passkey will be stronger than for the two-factor policy. For example, you can set the **User Verification Requirement** to **Required** when you configure the passwordless policy.
3. Configure the authentication flow. Use the **WebAuthn Browser** flow described in [Adding WebAuthn Authentication to a Browser Flow](#_webauthn-authenticator-setup). Configure the flow as follows:
   
   - The **WebAuthn Browser Forms** subflow contains **Username Form** as the first authenticator. Delete the default **Username Password Form** authenticator and add the **Username Form** authenticator. This action requires the user to provide a username as the first step.
   - There will be a required subflow, which can be named **Passwordless Or Two-factor**, for example. This subflow indicates the user can authenticate with Passwordless WebAuthn credential or with Two-factor authentication.
   - The flow contains **WebAuthn Passwordless Authenticator** as the first alternative.
   - The second alternative will be a subflow named **Password And Two-factor Webauthn**, for example. This subflow contains a **Password Form** and a **WebAuthn Authenticator**.

The final configuration of the flow looks similar to this:

PasswordLess flow

![PasswordLess flow](./images/webauthn-passwordless-flow.png)

You can now add **WebAuthn Register Passwordless** as the required action to a user, already known to Keycloak, to test this. During the first authentication, the user must use the password and second-factor WebAuthn credential. The user does not need to provide the password and second-factor WebAuthn credential if they use the WebAuthn Passwordless credential.

#### [](#_webauthn_loginless)LoginLess WebAuthn

Keycloak uses WebAuthn for two-factor authentication, but you can use WebAuthn as the first-factor authentication. In this case, users with `passwordless` WebAuthn credentials can authenticate to Keycloak without submitting a login or a password. Keycloak can use WebAuthn as both the loginless/passwordless and two-factor authentication mechanism in the context of a realm.

An administrator typically requires that Passkeys registered by users for the WebAuthn loginless authentication meet different requirements. Loginless authentication requires users to authenticate to the Passkey (for example by using a PIN code or a fingerprint) and that the cryptographic keys associated with the loginless credential are stored physically on the Passkey. Not all Passkeys meet that kind of requirement. Check with your Passkey vendor if your device supports 'user verification' and 'discoverable credential'. See [Supported Passkeys](#_webauthn-supported-keys).

Keycloak permits administrators to configure the `WebAuthn Passwordless Policy` in a way that allows loginless authentication. Note that loginless authentication can only be configured with `WebAuthn Passwordless Policy` and with `WebAuthn Passwordless` credentials. WebAuthn loginless authentication and WebAuthn passwordless authentication can be configured on the same realm but will share the same policy `WebAuthn Passwordless Policy`.

##### [](#setup-3)Setup

Procedure

Set up WebAuthn Loginless support as follows:

1. (If not already done) Check the required action for **WebAuthn Register Passwordless** is enabled. Use the steps described in [Enable WebAuthn Authenticator Registration](#_webauthn-register), but using **WebAuthn Register Passwordless** instead of **WebAuthn Register**.
2. Configure the `WebAuthn Passwordless Policy` if needed. Perform the configuration in the Admin Console, `Authentication` section, in the `Policies` → `WebAuthn Passwordless Policy` tab. By default, Keycloak sets **User Verification Requirement** to **required** and **Require Discoverable Credential** to **Yes** for the passwordless scenario to work properly. Storage capacity is usually very limited on Passkeys meaning that you won’t be able to store many discoverable credentials on your Passkey.
3. Configure the authentication flow. Create a new authentication flow, add the "WebAuthn Passwordless" execution and set the Requirement setting of the execution to **Required**

The final configuration of the flow looks similar to this:

LoginLess flow

![LoginLess flow](./images/webauthn-loginless-flow.png)

You can now add the required action `WebAuthn Register Passwordless` to a user, already known to Keycloak, to test this. The user with the required action configured will have to authenticate (with a username/password for example) and will then be prompted to register a Passkey to be used for loginless authentication.

##### [](#vendor-specific-remarks)Vendor specific remarks

###### [](#compatibility-check-list)Compatibility check list

Loginless authentication with Keycloak requires the Passkey to meet the following features

- FIDO2 compliance: not to be confused with FIDO/U2F
- User verification: the ability for the Passkey to authenticate the user (prevents someone finding your Passkey to be able to authenticate loginless and passwordless)
- Discoverable Credential: the ability for the Passkey to store the login and the cryptographic keys associated with the client application

###### [](#windows-hello)Windows Hello

To use Windows Hello based credentials to authenticate against Keycloak, configure the **Signature Algorithms** setting of the `WebAuthn Passwordless Policy` to include the **RS256** value. Note that some browsers don’t allow access to platform Passkey (like Windows Hello) inside private windows.

###### [](#_webauthn-supported-keys)Supported Passkeys

The following Passkeys have been successfully tested for loginless authentication with Keycloak:

- Windows Hello (Windows 10 21H1/21H2)
- Yubico Yubikey 5 NFC
- Feitian ePass FIDO-NFC

### [](#passkeys_server_administration_guide)Passkeys

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/passkeys.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fpasskeys.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fpasskeys.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak provides support for [Passkeys](https://fidoalliance.org/passkeys/). Keycloak works as a Passkeys Relying Party (RP).

Passkey registration and authentication are performed using the same features of [WebAuthn](#webauthn_server_administration_guide). More specifically **Passkeys** are related to [LoginLess WebAuthn](#_webauthn_loginless) as they try to avoid any password during login. Therefore, users of Keycloak can do Passkey registration and authentication by existing [WebAuthn registration and authentication](#webauthn_server_administration_guide), using the **passwordless** variants.

The **Passkeys** feature has been integrated seamlessly in the default authentication forms in two different ways. When activated, both conditional UI and modal UI are available in the forms in which the username input is displayed (for example **Username Password Form** or **Username Form**). Besides, the password forms, when the username was already selected, always show the modal UI button to login by passkey if the current user has passwordless WebAuthn credentials associated. This way modal and conditional UI can be used to perform a complete login from scratch that needs username and password, and only modal UI is presented when the username is already selected in the authentication process (because of re-authentication or because the user was selected before in the process not using a passkey).

**Passkeys** have been added to the following authenticator implementations:

1. **Username Password Form**: The username and password form used by default in Keycloak.
2. **Username Form**: The form in which the username is displayed alone and is typically followed by the password form. This authenticator is used when the username and password fields want to be presented to the user in two different steps.
3. **Password Form**: Authenticating using **Passkeys** in the **Username Form** skips the next **Password Form** execution. The **Password Form** implementation checks if the user was already authenticated using a passwordless WebAuthn credential and, if that is the case, no password is requested. If the **Password Form** cannot be skipped, it allows using modal UI to authenticate the user if the account has passkey credentials associated.
4. **Organization Identity - First Login**: The organization form that is used when the [organizations](#_enabling_organization_) feature is enabled for the realm. Using **Passkeys** in this step avoids the subsequent execution of the username and password form in the same way than in the username form.
5. **Username Password Form for identity provider re-authentication**: Similar to the default **Username Password Form** but used in the first login flow to re-authenticate and prove ownership of the account. Now the modal UI button is available to use passkeys to re-authenticate.

Finally, the default **browser** flow is modified to skip the flow **Browser - Conditional 2FA** if a passkey was used previously to authenticate the account. So now, when passkeys are enabled in the realm, 2FA is only presented if the user introduced a common password to login, not if a passkey was used. The new **Condition - credential** is added to the sub-flow to check the credential presented before. If you want to always use 2FA, you can just change the requirement from **Required** to **Disabled** for this condition. See [Conditions in conditional flows](#conditions-in-conditional-flows) for more information.

Both synced Passkeys and device-bound Passkeys can be used for both Same-Device and Cross-Device Authentication (CDA). However, Passkeys operations success depends on the user’s environment. Make sure which operations can succeed in [the environment](https://passkeys.dev/device-support/).

#### [](#_passkeys-conditional-ui)Passkey Authentication with Conditional UI or autofill

The Conditional User Interface (UI) or autofill is a feature related to passkeys in which the username input (the field in which the username to login is typed) is tagged with a `webauthn` autofill detail token (for example using the attribute `autocomplete="username webauthn"`). When the user clicks in such an input field, the user agent (browser) can render a list of discovered credentials for the user to select from, and perhaps also give the user the option to *try another way*. If the user selects one of the presented passkeys, Keycloak initiates the WebAuthn authentication with that key and avoids any password typing.

Compared with [LoginLess WebAuthn](#_webauthn_loginless), the authentication improves the user’s experience of authentication.

Passkey Authentication with Conditional UI Autofill using Chrome browser

![Passkey Authentication with Conditional UI Autofill using Chrome browser](./images/passkey-conditional-ui-autofill.png)

This authentication uses the [WebAuthn Conditional UI](https://github.com/w3c/webauthn/wiki/Explainer:-WebAuthn-Conditional-UI/). Therefore, this authentication success depends on the user’s environment. If the environment does not support WebAuthn Conditional UI, the user should use the direct modal UI or username and password login.

#### [](#passkeys-authentication-with-modal-ui)Passkeys Authentication with Modal UI

Nevertheless, because conditional UI can sometimes not show all the credentials to the user, the modal UI can always be initiated using the button **Sign in with Passkey**. The Modal User Interface (UI) ensures all passkeys are usable, including the ones stored in hardware tokens or on other devices that cannot be enumerated without user interaction.

The modal UI button is also presented in the password forms when the user is already selected, for example when re-authenticating in the common **Username Password Form**. In this case, the modal UI is limited to the passwordless WebAuthn credentials that Keycloak has defined for the account.

Passkey Authentication with Modal UI using Chrome browser

![Passkey Authentication with Modal UI using Chrome browser](./images/passkey-modal-ui.png)

#### [](#setup-4)Setup

Set up Passkey Authentication for the default forms as follows:

1. (If not already done) Check the required action for **WebAuthn Register Passwordless** is enabled. Use the steps described in [Enable WebAuthn Authenticator Registration](#_webauthn-register), but using **WebAuthn Register Passwordless** instead of **WebAuthn Register**.
2. Configure the **WebAuthn Passwordless Policy** in the same way that is explained in [LoginLess WebAuthn](#_webauthn_loginless). Perform the configuration in the Admin Console, `Authentication` section, in the tab `Policies` → `WebAuthn Passwordless Policy`. The default configuration for the passwordless policy is usually enough for correct passkeys integration.
   
   Storage capacity is usually very limited on hardware passkeys meaning that you cannot store many discoverable credentials on your passkey. However, this limitation may be mitigated for instance if you use an Android phone backed by a Google account as a passkey device or an iPhone backed by Bitwarden.
3. In the **WebAuthn Passwordless Policy** tab, activate the **Enable Passkeys** option at the bottom. This switch is the one that really enables passkeys in the realm.

### [](#_recovery-codes)Recovery Codes

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/recovery-codes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Frecovery-codes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Frecovery-codes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The Recovery Codes are a number of sequential one-time passwords (currently 12) auto-generated by Keycloak. The codes can be used as a 2nd Factor Authentication (2FA) by adding the `Recovery Authentication Code Form` authenticator to your authentication flow. When configured in the flow, Keycloak asks the user for the next generated code in order. When the current code is introduced by the user, it is removed and the next code will be required for the next login.

Due to its nature, the Recovery Codes work normally as a backup for another 2FA methods. They can complement the `OTP Form` or the `WebAuthn Authenticator` to give a backing way to log inside Keycloak, for example, if the software or hardware device used for the previous 2FA methods is broken or unavailable.

You can configure the `Configure OTP` required action to ask for recovery codes automatically by enabling **Add recovery codes**.

#### [](#check-recovery-codes-required-action-is-enabled)Check Recovery Codes required action is enabled

Check the Recovery Codes action is enabled in Keycloak:

1. Click **Authentication** in the menu.
2. Click the **Required Actions** tab.
3. Ensure the **Recovery Authentication Codes** switch **Enabled** is set to **On**.

Toggle the **Default Action** switch to **On** if you want all the new users to register their Recovery Codes credentials in the first login.

#### [](#configure-the-recovery-codes-required-action)Configure the Recovery Codes required action

From the **Required Actions** tab of the admin console, you have the option to configure the **Recovery Authentication Codes** required action. So far, there is a configuration option **Warning Threshold** available. When a user has a smaller amount of remaining recovery codes on his account than the value configured here, account console will show warning to the user, which will recommend them to set up a new set of recovery codes. The warning displayed to the user may look similar to this:

Recovery Codes Account console warning

![Recovery Codes Account console warning](./images/recovery-codes-account-console-warn.png)

#### [](#adding-recovery-codes-to-the-browser-flow)Adding Recovery Codes to the browser flow

The following procedure adds the `Recovery Authentication Code Form` as an alternative way of login in the default **Browser** flow.

1. Click **Authentication** in the realm menu.
2. Click the **Browser** flow.
3. Locate the execution **Recovery Authentication Code Form** inside the **Browser - Conditional 2FA** sub-flow.
4. Change the *requirement* from *Disabled* to *Alternative* for that execution.
   
   Recovery Codes Browser flow
   
   ![Recovery Codes Browser flow](./images/recovery-codes-browser-flow.png)
   
   With this configuration, both 2FA authenticators (`OTP Form` and `Recovery Authentication Code Form`) are alternate ways to log into Keycloak. If the user has configured both credential types, the credential with the highest priority will be displayed by default, but the **Try Another Way** option will appear so that the user has the alternative methods to log in.

You can see more examples of 2FA configurations in [2FA conditional workflow examples](#twofa-conditional-workflow-examples).

#### [](#creating-the-recovery-codes-credential)Creating the Recovery Codes credential

Once the Recovery Codes required action is enabled and the credential type is managed in the flow, users can request to create their own codes. The action is just another [required action](#con-required-actions_server_administration_guide) that can be used in Keycloak (directly called by the user by using the Account Console or assigned by an administrator by using the Admin Console).

The required action, when executed, generates the list of codes and presents it to the user. The action offers to print, download, or copy the list of codes to help the user to store them is a safe place. In order to complete the setup, the checkbox **I have saved these codes somewhere safe** should be previously checked.

Recovery Authentication Codes setup page

![Recovery Authentication Codes setup page](./images/recovery-codes-setup.png)

The Recovery Codes can be re-created at any moment.

### [](#conditions-in-conditional-flows)Conditions in conditional flows

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/conditions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fconditions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fconditions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

As was mentioned in [Execution requirements](#_execution-requirements), *Condition* executions can be only contained in *Conditional* subflow. If all *Condition* executions evaluate as true, then the *Conditional* sub-flow acts as *Required*. You can process the next execution in the *Conditional* sub-flow. If some executions included in the *Conditional* sub-flow evaluate as false, then the whole sub-flow is considered as *Disabled*.

#### [](#available-conditions)Available conditions

`Condition - User Role`

This execution has the ability to determine if the user has a role defined by *User role* field. If the user has the required role, the execution is considered as true and other executions are evaluated. The administrator has to define the following fields:

Alias

Describes a name of the execution, which will be shown in the authentication flow.

User role

Role the user should have to execute this flow. To specify an application role the syntax is `appname.approle` (for example `myapp.myrole`).

`Condition - User Configured`

This checks if the other executions in the flow are configured for the user. The Execution requirements section includes an example of the OTP form.

`Condition - User Attribute`

This checks if the user has set up the required attribute: optionally, the check can also evaluate the group attributes. There is a possibility to negate output, which means the user should not have the attribute. The [User Attributes](#user-profile) section shows how to add a custom attribute. You can provide these fields:

Alias

Describes a name of the execution, which will be shown in the authentication flow.

Attribute name

Name of the attribute to check.

Expected attribute value

Expected value in the attribute.

Include group attributes

If On, the condition checks if any of the joined group has one attribute matching the configured name and value: this option can affect performance

Negate output

You can negate the output. In other words, the attribute should not be present.

`Condition - sub-flow executed`

The condition checks if a previous sub-flow was successfully executed (or not executed) in the authentication process. Therefore, the flow can trigger other steps based on a previous sub-flow termination. These configuration fields exist:

Flow name

The sub-flow name to check if it was executed or not executed. Required.

Check result

When the condition evaluates to true. If `executed` returns true when the configured sub-flow was executed with output success, false otherwise. If `not-executed` returns false when the sub-flow was executed with output success, true otherwise (the negation of the previous option).

`Condition - client scope`

The condition to evaluate if a configured client scope is present as a client scope of the client requesting authentication. These configuration fields exist:

Client scope name

The name of the client scope, which should be present as a client scope of the client, which is requesting authentication. If requested client scope is default client scope of the client requesting login, the condition will be evaluated to true. If requested client scope is optional client scope of the client requesting login, condition will be evaluated to true if client scope is sent by the client in the login request (for example, by the `scope` parameter in case of OIDC/OAuth2 client login). Required.

Negate output

Apply a NOT to the check result. When this is true, then the condition will evaluate to true just if configured client scope is not present.

`Condition - credential`

This condition evaluates if a specific credential type has been used (or not used) by the user during the authentication process. Configuration options:

Credentials

The list of credentials to be considered by the condition.

Included

If **included** is true, the condition will be evaluated to `true` when any of the credentials specified in the **credentials** option has been used in the authentication process, false otherwise. If **included** is false, the condition is evaluated in the opposite way, it will be `true` if none of the **credentials** configured have been used, and `false` if one or more of them have been used.

#### [](#explicitly-denyallow-access-in-conditional-flows)Explicitly deny/allow access in conditional flows

You can allow or deny access to resources in a conditional flow. The two authenticators `Deny Access` and `Allow Access` control access to the resources by conditions.

`Allow Access`

Authenticator will always successfully authenticate. This authenticator is not configurable.

`Deny Access`

Access will always be denied. You can define an error message, which will be shown to the user. You can provide these fields:

Alias

Describes a name of the execution, which will be shown in the authentication flow.

Error message

Error message which will be shown to the user. The error message could be provided as a particular message or as a property in order to use it with localization. (i.e. "*You do not have the role 'admin'.*", *my-property-deny* in messages properties) Leave blank for the default message defined as property `access-denied`.

Here is an example how to deny access to all users who do not have the role `role1` and show an error message defined by a property `deny-role1`. This example includes `Condition - User Role` and `Deny Access` executions.

Browser flow

![Deny access flow](./images/deny-access-flow.png)

Condition - user role configuration

![Deny access role settings](./images/deny-access-role-condition.png)

Configuration of the `Deny Access` is really easy. You can specify an arbitrary Alias and required message like this:

![Deny access execution settings](./images/deny-access-execution-cond.png)

The last thing is defining the property with an error message in the login theme `messages_en.properties` (for English):

```
deny-role1 = You do not have required role!
```

#### [](#twofa-conditional-workflow-examples)2FA conditional workflow examples

The section presents some examples of conditional workflows that integrates 2nd Factor Authentication (2FA) in different ways. The examples copy the default `browser` flow and modify the configuration inside the `forms` sub-flow.

##### [](#conditional-2fa-sub-flow)Conditional 2FA sub-flow

The default `browser` flow uses a `Conditional 2FA` sub-flow that already gives 2nd factor Authentication (2FA) with OTP Form (One Time Password). It also provides WebAuthn and Recovery Codes but they are disabled by default. Consistent with this approach, different 2FA methods can be integrated with the `Condition - User Configured`.

2FA all alternative

![2FA all alternative](./images/2fa-example1.png)

The `forms` sub-flow contains another `2FA` conditional sub-flow with `Condition - user configured`. Three 2FA steps (OTP, Webauthn and Recovery Codes) are allowed as alternative steps. The user will be able to choose one of the three options, if they are configured for the user. As the sub-flow is conditional, the authentication process will complete successfully if no 2FA credential is configured.

This configuration provides the same behavior as when you configure with the default **browser** flow with both *Disabled* steps are configured to *Alternative*.

##### [](#conditional-2fa-sub-flow-and-deny-access)Conditional 2FA sub-flow and deny access

The second example continues the previous one. After the `2FA` sub-flow, another flow `Deny access if no 2FA` is used to check if the previous `2FA` was not executed. In that case (the user has no 2FA credential configured) the access is denied.

2FA all alternative and deny access

![2FA all alternative and deny access](./images/2fa-example2.png)

The `Condition - sub-flow executed` is configured to detect if the `2FA` sub-flow was not executed previously.

Configuration for the sub-flow executed

![Configuration for the sub-flow executed](./images/2fa-example2-config.png)

The step `Deny access` denies the authentication if not executed.

##### [](#_conditional-2fa-otp-default)Conditional 2FA sub-flow with OTP default

The last example is very similar to the previous one. Instead of denying the access, step `OTP Form` is configured as required.

2FA all alternative with OTP default

![2FA all alternative with OTP default](./images/2fa-example3.png)

With this flow, if the user has none of the 2FA methods configured, the OTP setup will be enforced to continue the login.

### [](#_authentication-sessions)Authentication sessions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/authentication/authentication-sessions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fauthentication%2Fauthentication-sessions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fauthentication%2Fauthentication-sessions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When a login page is opened for the first time in a web browser, Keycloak creates an object called authentication session that stores some useful information about the request. Whenever a new login page is opened from a different tab in the same browser, Keycloak creates a new record called authentication sub-session that is stored within the authentication session. Authentication requests can come from any type of clients such as the Admin CLI. In that case, a new authentication session is also created with one authentication sub-session. Please note that authentication sessions can be created also in other ways than using a browser flow.

The authentication session usually expires after 30 minutes by default. The exact time is specified by the **Login timeout** switch in the **Sessions** tab of the admin console when [configuring realms](#_configuring-realms).

#### [](#authentication-in-more-browser-tabs)Authentication in more browser tabs

As described in the previous section, a situation can involve a user who is trying to authenticate to the Keycloak server from multiple tabs of a single browser. However, when that user authenticates in one browser tab, the other browser tabs will automatically restart the authentication. This authentication occurs due to the small javascript available on the Keycloak login pages. The restart will typically authenticate the user in other browser tabs and redirect to clients because there is an SSO session now due to the fact that the user just successfully authenticated in first browser tab. Some rare exceptions exist when a user is not automatically authenticated in other browser tabs, such as for instance when using an OIDC parameter *prompt=login* or [step-up authentication](#_step-up-flow) requesting a stronger authentication factor than the currently authenticated factor.

In some rare cases, it can happen that after authentication in the first browser tab, other browser tabs are not able to restart authentication because the authentication session is already expired. In this case, the particular browser tab will redirect the error about the expired authentication session back to the client in a protocol specific way. For more details, see the corresponding sections of **OIDC documentation** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section. When the client application receives such an error, it can immediately resubmit the OIDC/SAML authentication request to Keycloak as this should usually automatically authenticate the user due to the existing SSO session as described earlier. As a result, the end user is authenticated automatically in all browser tabs. The **Keycloak JavaScript adapter** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section, and [Keycloak Identity provider](#_identity_broker) support to handle this error automatically and retry the authentication to the Keycloak server in such a case.

## [](#_identity_broker)Integrating identity providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

An Identity Broker is an intermediary service connecting service providers with identity providers. The identity broker creates a relationship with an external identity provider to use the provider’s identities to access the internal services the service provider exposes.

From a user perspective, identity brokers provide a user-centric, centralized way to manage identities for security domains and realms. You can link an account with one or more identities from identity providers or create an account based on the identity information from them.

An identity provider derives from a specific protocol used to authenticate and send authentication and authorization information to users. It can be:

- A social provider such as Facebook, Google, or Twitter.
- A business partner whose users need to access your services.
- A cloud-based identity service you want to integrate.

Typically, Keycloak bases identity providers on the following protocols:

- `SAML v2.0`
- `OpenID Connect v1.0`
- `OAuth v2.0`

### [](#_identity_broker_overview)Brokering overview

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/overview.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Foverview.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Foverview.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When using Keycloak as an identity broker, Keycloak does not force users to provide their credentials to authenticate in a specific realm. Keycloak displays a list of identity providers from which they can authenticate.

If you configure a default identity provider, Keycloak redirects users to the default provider.

Different protocols may require different authentication flows. All the identity providers supported by Keycloak use the following flow.

Identity broker flow

![Identity broker flow](./images/identity_broker_flow.png)

01. The unauthenticated user requests a protected resource in a client application.
02. The client application redirects the user to Keycloak to authenticate.
03. Keycloak displays the login page with a list of identity providers configured in a realm.
04. The user selects one of the identity providers by clicking its button or link.
05. Keycloak issues an authentication request to the target identity provider requesting authentication and redirects the user to the identity provider’s login page. The administrator has already set the connection properties and other configuration options for the Admin Console’s identity provider.
06. The user provides credentials or consents to authenticate with the identity provider.
07. Upon successful authentication by the identity provider, the user redirects back to Keycloak with an authentication response. Usually, the response contains a security token used by Keycloak to trust the identity provider’s authentication and retrieve user information.
08. Keycloak checks if the response from the identity provider is valid. If valid, Keycloak imports and creates a user if the user does not already exist. Keycloak may ask the identity provider for further user information if the token does not contain that information. This behavior is *identity federation*. If the user already exists, Keycloak may ask the user to link the identity returned from the identity provider with the existing account. This behavior is *account linking*. With Keycloak, you can configure *Account linking* and specify it in the [First Login Flow](#_identity_broker_first_login). At this step, Keycloak authenticates the user and issues its token to access the requested resource in the service provider.
09. When the user authenticates, Keycloak redirects the user to the service provider by sending the token previously issued during the local authentication.
10. The service provider receives the token from Keycloak and permits access to the protected resource.

Variations of this flow are possible. For example, the client application can request a specific identity provider rather than displaying a list of them, or you can set Keycloak to force users to provide additional information before federating their identity.

At the end of the authentication process, Keycloak issues its token to client applications. Client applications are separate from the external identity providers, so they cannot see the client application’s protocol or how they validate the user’s identity. The provider only needs to know about Keycloak.

### [](#default_identity_provider)Default Identity Provider

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/default-provider.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fdefault-provider.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fdefault-provider.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak can redirect to an identity provider rather than displaying the login form. To enable this redirection:

Procedure

1. Click **Authentication** in the menu.
2. Click the **Browser** flow.
3. Click the gear icon **⚙️** on the **Identity Provider Redirector** row.
4. Set **Default Identity Provider** to the identity provider you want to redirect users to.

If Keycloak does not find the configured default identity provider, the login form is displayed.

This authenticator is responsible for processing the `kc_idp_hint` query parameter. See the [client suggested identity provider](#_client_suggested_idp) section for more information.

The authenticator will redirect to the identity provider and authentication is delegated to the identity provider. The `browser` authentication flow will not continue after the login with the identity provider is successfully finished. If you want to perform additional steps after the identity provider login (for example 2-factor authentication), it may be needed to configure [Post login flow](#_identity_broker_post_login_flow).

### [](#_general-idp-config)General configuration

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/configuration.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fconfiguration.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fconfiguration.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The foundations of the identity broker configuration are identity providers (IDPs). Keycloak creates identity providers for each realm and enables them for every application by default. Users from a realm can use any of the registered identity providers when signing in to an application.

Procedure

1. Click **Identity Providers** in the menu.
   
   Identity Providers
   
   ![Identity Providers](./images/identity-providers.png)
2. Select an identity provider. Keycloak displays the configuration page for the identity provider you selected.
   
   Add Facebook identity Provider
   
   ![Add Facebook Identity Provider](./images/add-identity-provider.png)
   
   When you configure an identity provider, the identity provider appears on the Keycloak login page as an option. You can place custom icons on the login screen for each identity provider. See [custom icons](https://www.keycloak.org/ui-customization/themes#custom-identity-providers-icons) for more information.
   
   IDP login page
   
   ![identity provider login page](./images/identity-provider-login-page.png)
   
   Social
   
   Social providers enable social authentication in your realm. With Keycloak, users can log in to your application using a social network account. Supported providers include Twitter, Facebook, Google, LinkedIn, Instagram, Microsoft, PayPal, Openshift v4, GitHub, GitLab, Bitbucket, and Stack Overflow.
   
   Protocol-based
   
   Protocol-based providers rely on specific protocols to authenticate and authorize users. Using these providers, you can connect to any identity provider compliant with a specific protocol. Keycloak provides support for SAML v2.0 and OpenID Connect v1.0 protocols. You can configure and broker any identity provider based on these open standards.

Although each type of identity provider has its configuration options, all share a common configuration. The following configuration options available:

Table 1. Common Configuration   Configuration Description

Alias

The alias is a unique identifier for an identity provider and references an internal identity provider. Keycloak uses the alias to build redirect URIs for OpenID Connect protocols that require a redirect URI or callback URL to communicate with an identity provider. All identity providers must have an alias. Alias examples include `facebook`, `google`, and `idp.acme.com`.

Enabled

Toggles the provider ON or OFF.

Hide on Login Page

When **ON**, Keycloak does not display this provider as a login option on the login page. Clients can request this provider by using the 'kc\_idp\_hint' parameter in the URL to request a login.

Account Linking Only

When **ON**, Keycloak links existing accounts with this provider. This provider cannot log users in, and Keycloak does not display this provider as an option on the login page.

Store Tokens

When **ON**, Keycloak stores tokens from the identity provider.

Stored Tokens Readable

When **ON**, users can retrieve the stored identity provider token. This action also applies to the *broker* client-level role *read token*.

Trust Email

When **ON**, Keycloak trusts email addresses from the identity provider. If the realm requires email validation, users that log in from this identity provider do not need to perform the email verification process. If the target identity provider supports email verification and advertises this information when returning the user profile information, the email of the federated user will be (un)marked as verified. For instance, an OpenID Connect Provider returning a `email_verified` claim in their ID Tokens. Note that this setting will set the email as verified when the user is federated for the first time and on subsequent logins through the broker if the sync mode is set to `FORCE`.

GUI Order

The sort order of the available identity providers on the login page.

Verify essential claim

When **ON**, ID tokens issued by the identity provider must have a specific claim, otherwise, the user can not authenticate through this broker

Essential claim

When **Verify essential claim** is **ON**, the name of the JWT token claim to filter (match is case sensitive)

Essential claim value

When **Verify essential claim** is **ON**, the value of the JWT token claim to match (supports regular expression format)

First Login Flow

The authentication flow Keycloak triggers when users use this identity provider to log into Keycloak for the first time.

Post Login Flow

The authentication flow Keycloak triggers when a user finishes logging in with the external identity provider.

Sync Mode

Strategy to update user information from the identity provider through mappers. When choosing **legacy**, Keycloak used the current behavior. **Import** does not update user data and **force** updates user data when possible. See [Identity Provider Mappers](#_mappers) for more information.

Case-sensitive username

If enabled, the original username from the identity provider is kept as is when federating users. Otherwise, the username from the identity provider is lower-cased and might not match the original value if it is case-sensitive. This setting only affects the username associated with the federated identity as usernames in the server are always in lower-case.

Show in Account console

Defines how the identity provider will be available from the account console. If set to `Always`, it is always available. If set to `When linked`, it will be shown only when the user is linked to it. Otherwise, if set to `Never`, it will not be available even if the user is linked to it.

### [](#social-identity-providers)Social Identity Providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social-login.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial-login.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial-login.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

A social identity provider can delegate authentication to a trusted, respected social media account. Keycloak includes support for social networks such as Google, Facebook, Twitter, GitHub, LinkedIn, Microsoft, and Stack Overflow.

#### [](#bitbucket)Bitbucket

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/bitbucket.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fbitbucket.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fbitbucket.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To log in with Bitbucket, perform the following procedure.

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **Bitbucket**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/bitbucket-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, perform the [OAuth on Bitbucket Cloud](https://support.atlassian.com/bitbucket-cloud/docs/use-oauth-on-bitbucket-cloud/) process. When you click **Add Consumer**:
   
   1. Paste the value of **Redirect URI** into the **Callback URL** field.
   2. Ensure you select **Email** and **Read** in the **Account** section to permit your application to read email.
5. Note the **Key** and **Secret** values Bitbucket displays when you create your consumer.
6. In Keycloak, paste the value of the `Key` into the **Client ID** field.
7. In Keycloak, paste the value of the `Secret` into the **Client Secret** field.
8. Click **Add**.

#### [](#_facebook)Facebook

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/facebook.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Ffacebook.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Ffacebook.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **Facebook**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/facebook-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, open the [Meta for Developers](https://developers.facebook.com/).
   
   01. Click **My Apps**.
   02. Select **Create App**.
       
       Add a use case
       
       ![Add a use case](./images/meta-add-use-case.png)
   03. Select **Other**.
       
       Select an app type
       
       ![Select an app type](./images/meta-select-app-type.png)
   04. Select **Consumer**.
       
       Create an app
       
       ![Create an app](./images/meta-create-app.png)
   05. Fill in all required fields.
   06. Click **Create app**. Meta then brings you to the dashboard.
       
       Add a product
       
       ![Add Product](./images/meta-add-product.png)
   07. Click **Set Up** in the **Facebook Login** box.
   08. Select **Web**.
   09. Enter the **Redirect URI’s** value into the **Site URL** field and click **Save**.
   10. In the navigation panel, select **App settings** - **Basic**.
   11. Click **Show** in the **App Secret** field.
   12. Note the **App ID** and the **App Secret**.
5. Enter the [`App ID` and `App Secret`](https://developers.facebook.com/documentation/facebook-login/guides/access-tokens) values from your Facebook app into the **Client ID** and **Client Secret** fields in Keycloak.
6. Click **Add**
7. Enter the required scopes into the **Default Scopes** field. By default, Keycloak uses the **email** scope. See [Graph API](https://developers.facebook.com/docs/graph-api) for more information about Facebook scopes.

Keycloak sends profile requests to `graph.facebook.com/me?fields=id,name,email,first_name,last_name` by default. The response contains the id, name, email, first\_name, and last\_name fields only. To fetch additional fields from the Facebook profile, add a corresponding scope and add the field name in the `Additional user’s profile fields` configuration option field.

#### [](#_github)GitHub

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/github.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fgithub.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fgithub.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To log in with GitHub, perform the following procedure.

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **Github**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/github-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, [create an OAuth app](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app) or [create an GitHub app](https://docs.github.com/en/apps/creating-github-apps/about-creating-github-apps/about-creating-github-apps). Note that only GitHub apps can refresh tokens, while OAuth apps cannot refresh tokens.
   
   1. Enter the value of **Redirect URI** into the **Authorization callback URL** field when creating the app.
   2. Note the **Client ID** and **Client secret** on the management page of your OAUTH app.
5. In Keycloak, paste the value of the `Client ID` into the **Client ID** field.
6. In Keycloak, paste the value of the `Client secret` into the **Client Secret** field.
7. Enable **JSON Format** to retrieve the external IDP Tokens in JSON format. Note that this is also required for tokens to be refreshed.
8. Click **Add**.

#### [](#gitlab)GitLab

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/gitlab.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fgitlab.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fgitlab.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **GitLab**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/gitlab-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, [add a new GitLab application](https://docs.gitlab.com/integration/oauth_provider/).
   
   1. Use the **Redirect URI** in your clipboard as the **Redirect URI**.
   2. Note the **Application ID** and **Secret** when you save the application.
5. In Keycloak, paste the value of the `Application ID` into the **Client ID** field.
6. In Keycloak, paste the value of the `Secret` into the **Client Secret** field.
7. Click **Add**.

#### [](#_google)Google

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/google.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fgoogle.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fgoogle.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Procedure

01. Click **Identity Providers** in the menu.
02. From the **Add provider** list, select **Google**.
    
    Add identity provider
    
    ![Add Identity Provider](./images/google-add-identity-provider.png)
03. Copy the value of **Redirect URI** to your clipboard.
04. In a separate browser tab open [the Google Cloud Platform console](https://console.cloud.google.com/).
05. In the Google dashboard for your Google app, in the Navigation menu on the left side, hover over **APIs & Services** and then click on the **OAuth consent screen** option. Create a consent screen, ensuring that the user type of the consent screen is **External**.
06. In the Google dashboard:
    
    1. Click the **Credentials** menu.
    2. Click **CREATE CREDENTIALS** - **OAuth Client ID**.
    3. From the **Application type** list, select **Web application**.
    4. Use the **Redirect URI** in your clipboard as the **Authorized redirect URIs**
    5. Click **Create**.
    6. Note **Your Client ID** and **Your Client secret**.
07. In Keycloak, paste the value of the `Your Client ID` into the **Client ID** field.
08. In Keycloak, paste the value of the `Your Client secret` into the **Client Secret** field.
09. Click **Add**
10. Enter the required scopes into the **Default Scopes** field. By default, Keycloak uses the following scopes: **openid** **profile** **email**. See the [OAuth Playground](https://developers.google.com/oauthplayground/) for a list of Google scopes.
11. To restrict access to your GSuite organization’s members only, enter the G Suite domain into the **Hosted Domain** field.
12. Click **Save**.

#### [](#instagram)Instagram

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/instagram.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Finstagram.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Finstagram.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The Instagram Identity Broker is deprecated for removal. Prefer using the Facebook Identity Broker instead. To enable it, start the server with `--features=instagram-broker`.

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **Instagram**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/instagram-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, open the [Meta for Developers](https://developers.facebook.com/).
   
   01. Click **My Apps**.
   02. Select **Create App**.
       
       Add a use case
       
       ![Add a use case](./images/meta-add-use-case.png)
   03. Select **Other**.
       
       Select an app type
       
       ![Select an app type](./images/meta-select-app-type.png)
   04. Select **Consumer**.
       
       Create an app
       
       ![Create an app](./images/meta-create-app.png)
   05. Fill in all required fields.
   06. Click **Create app**. Meta then brings you to the dashboard.
   07. In the navigation panel, select **App settings** - **Basic**.
   08. Select **+ Add Platform** at the bottom of the page.
   09. Click **\[Website]**.
   10. Enter a URL for your site.
       
       Add a product
       
       ![Add Product](./images/meta-add-product.png)
   11. Select **Dashboard** from the menu.
   12. Click **Set Up** in the **Instagram Basic Display** box.
   13. Click **Create New App**.
       
       Create a New Instagram App ID
       
       ![Create a New Instagram App ID](./images/instagram-create-instagram-app-id.png)
   14. Enter a value into the **Display name** field.
       
       Set up the app
       
       ![Setup the App](./images/instagram-app-settings.png)
   15. Paste the **Redirect URL** from Keycloak into the **Valid OAuth Redirect URIs** field.
   16. Paste the **Redirect URL** from Keycloak into the **Deauthorize Callback URL** field.
   17. Paste the **Redirect URL** from Keycloak into the **Data Deletion Request URL** field.
   18. Click **Show** in the **Instagram App Secret** field.
   19. Note the **Instagram App ID** and the **Instagram App Secret**.
   20. Click **App Review** - **Requests**.
   21. Follow the instructions on the screen.
5. In Keycloak, paste the value of the `Instagram App ID` into the **Client ID** field.
6. In Keycloak, paste the value of the `Instagram App Secret` into the **Client Secret** field.
7. Click **Add**.

#### [](#_linkedin)LinkedIn

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/linked-in.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Flinked-in.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Flinked-in.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **LinkedIn**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/linked-in-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, [create an app](https://developer.linkedin.com) in the LinkedIn developer portal.
   
   1. After you create the app, click the **Auth** tab.
   2. Enter the value of **Redirect URI** into the **Authorized redirect URLs for your app** field.
   3. Note **Your Client ID** and **Your Client Secret**.
   4. Click the **Products** tab and **Request access** for the **Sign In with LinkedIn using OpenID Connect** product.
5. In Keycloak, paste the value of the `Client ID` into the **Client ID** field.
6. In Keycloak, paste the value of the `Client Secret` into the **Client Secret** field.
7. Click **Add**.

#### [](#_microsoft)Microsoft

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/microsoft.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fmicrosoft.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fmicrosoft.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **Microsoft**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/microsoft-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, register an app on [Microsoft Azure](https://azure.microsoft.com/en-us) under **App registrations**.
   
   1. In the Redirect URI section, select **Web** as a platform and paste the value of **Redirect URI** into the field.
   2. Find you application under **App registrations** and add a new client secret in the **Certificates & secrets** section.
   3. Note the **Value** of the created secret.
   4. Note the **Application (client) ID** in the **Overview** section.
5. In Keycloak, paste the value of the `Application (client) ID` into the **Client ID** field.
6. In Keycloak, paste the `Value` of the secret into the **Client Secret** field.
7. Click **Add**.

#### [](#openshift-4)OpenShift 4

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/openshift.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fopenshift.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fopenshift.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Prerequisites

1. A certificate of the OpenShift 4 instance stored in the Keycloak Truststore.
2. A Keycloak server configured in order to use the truststore. For more information, see the [Configuring a Truststore](https://www.keycloak.org/server/keycloak-truststore) guide.

Procedure

1. Locate the Openshift 4 instance’s API URL by using this command:
   
   ```
   oc cluster-info
   ```
2. Look for the URL in a line that has this format:
   
   ```
   Kubernetes master is running at https://api.<your-openshift-domain>:6443
   ```
3. In the Admin Console, click **Identity Providers** in the menu.
4. From the **Add provider** list, select **Openshift v4**.
5. Enter the **Client ID** and **Client Secret** and in the **Base URL** field, enter the API URL of your OpenShift 4 instance. Additionally, you can copy the **Redirect URI** to your clipboard.
   
   Add identity provider
   
   ![Add Identity Provider](./images/openshift-4-add-identity-provider.png)
6. Register your client, either via OpenShift 4 Console (Home → API Explorer → OAuth Client → Instances) or using the `oc` command-line tool.
   
   ```
   $ oc create -f <(echo '
   kind: OAuthClient
   apiVersion: oauth.openshift.io/v1
   metadata:
    name: kc-client (1)
   secret: "..." (2)
   redirectURIs:
    - "<here you can paste the Redirect URI that you copied in the previous step>" (3)
   grantMethod: prompt (4)
   ')
   ```

*1* The `name` of your OAuth client. Passed as `client_id` request parameter when making requests to `<openshift_master>/oauth/authorize` and `<openshift_master>/oauth/token`. The `name` parameter must be the same in the `OAuthClient` object and the Keycloak configuration. *2* The `secret` Keycloak uses as the `client_secret` request parameter. *3* The `redirect_uri` parameter specified in requests to `<openshift_master>/oauth/authorize` and `<openshift_master>/oauth/token` must be equal to (or prefixed by) one of the URIs in `redirectURIs`. The easiest way to configure it correctly is to copy-paste it from Keycloak OpenShift 4 Identity Provider configuration page (`Redirect URI` field). *4* The `grantMethod` Keycloak uses to determine the action when this client requests tokens but has not been granted access by the user.

In the end you should see the OpenShift 4 Identity Provider on the login page of your Keycloak instance. After clicking on it, you should be redirected to the OpenShift 4 login page.

Result

![Result](./images/openshift-4-result.png)

See [official OpenShift documentation](https://docs.okd.io/latest/authentication/configuring-oauth-clients.html#oauth-register-additional-client_configuring-oauth-clients) for more information.

#### [](#paypal)PayPal

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/paypal.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fpaypal.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fpaypal.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **PayPal**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/paypal-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, open the [PayPal Developer applications area](https://developer.paypal.com/developer/applications).
   
   1. Click **Create App** to create a PayPal app.
   2. Note the **Client ID** and **Client Secret**. Click the **Show** link to view the secret.
   3. Ensure **Log in with PayPal** is checked.
   4. Under Log in with PayPal click on **Advanced Settings**.
   5. Set the value of the **Return URL** field to the value of **Redirect URI** from Keycloak. Note that the URL can not contain `localhost`. If you want to use Keycloak locally, replace the `localhost` in the **Return URL** by `127.0.0.1` and then access Keycloak using `127.0.0.1` in the browser instead of `localhost`.
   6. Ensure **Full Name** and **Email** fields are checked.
   7. Click **Save** and then **Save Changes**.
5. In Keycloak, paste the value of the `Client ID` into the **Client ID** field.
6. In Keycloak, paste the value of the `Secret key 1` into the **Client Secret** field.
7. Click **Add**.

#### [](#_stackoverflow)Stack Overflow

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/stack-overflow.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fstack-overflow.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Fstack-overflow.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **Stack Overflow**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/stack-overflow-add-identity-provider.png)
3. In a separate browser tab, log into [registering your application on Stack Apps](https://stackapps.com/apps/oauth/register).
   
   Register application
   
   ![Register Application](./images/stack-overflow-app-register.png)
   
   1. Enter your application name into the **Application Name** field.
   2. Enter the OAuth domain into the **OAuth Domain** field.
   3. Click **Register Your Application**.
      
      Settings
      
      ![Settings](./images/stack-overflow-app-settings.png)
4. Note the **Client Id**, **Client Secret**, and **Key**.
5. In Keycloak, paste the value of the `Client Id` into the **Client ID** field.
6. In Keycloak, paste the value of the `Client Secret` into the **Client Secret** field.
7. In Keycloak, paste the value of the `Key` into the **Key** field.
8. Click **Add**.

#### [](#_twitter)Twitter

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/social/twitter.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Ftwitter.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsocial%2Ftwitter.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Prerequisites

1. A Twitter developer account.

Procedure

1. Click **Identity Providers** in the menu.
2. From the **Add provider** list, select **Twitter**.
   
   Add identity provider
   
   ![Add Identity Provider](./images/twitter-add-identity-provider.png)
3. Copy the value of **Redirect URI** to your clipboard.
4. In a separate browser tab, create an app in [Twitter Application Management](https://developer.twitter.com/apps/).
   
   1. Enter App name and click **Next**.
   2. Note the value of **API Key** and **API Key Secret** and click **App settings**.
   3. In the **User authentication settings** section click on the **Set up** button.
   4. Select **Web App** as the **Type of App**.
   5. Paste the value of the **Redirect URL** into the **Callback URI / Redirect URL** field.
   6. The value for **Website URL** can be any valid URL except `localhost`.
   7. Click **Save** and then **Done**.
5. In Keycloak, paste the value of the `API Key` into the **Client ID** field.
6. In Keycloak, paste the value of the `API Key Secret` into the **Client Secret** field.
7. Click **Add**.

### [](#_identity_broker_oidc)OpenID Connect v1.0 identity providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/oidc.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Foidc.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Foidc.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak brokers identity providers based on the OpenID Connect protocol. These identity providers (IDPs) must support the [Authorization Code Flow](#con-oidc-auth-flows_server_administration_guide) defined in the specification to authenticate users and authorize access.

Procedure

1. Click **Identity Providers** in the menu.
2. From the `Add provider` list, select `OpenID Connect v1.0`.
   
   Add identity provider
   
   ![Add Identity Provider](./images/oidc-add-identity-provider.png)
3. Enter your initial configuration options. See [General IDP Configuration](#_general-idp-config) for more information about configuration options.
   
   Table 2. OpenID connect config   Configuration Description
   
   Authorization URL
   
   The authorization URL endpoint the OIDC protocol requires.
   
   Token URL
   
   The token URL endpoint the OIDC protocol requires.
   
   Logout URL
   
   The logout URL endpoint in the OIDC protocol. This value is optional.
   
   Backchannel Logout
   
   A background, out-of-band, REST request to the IDP to log out the user. Some IDPs perform logout through browser redirects only, as they may identify sessions using a browser cookie.
   
   User Info URL
   
   An endpoint the OIDC protocol defines. This endpoint points to user profile information.
   
   Client Authentication
   
   Defines the Client Authentication method Keycloak uses with the Authorization Code Flow. In the case of JWT signed with a private key, Keycloak uses the realm private key. In the other cases, define a client secret. See the [Client Authentication specifications](https://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication) for more information.
   
   Client ID
   
   A realm acting as an OIDC client to the external IDP. The realm must have an OIDC client ID if you use the Authorization Code Flow to interact with the external IDP.
   
   Client Secret
   
   Client secret from an external [vault](#_vault-administration). This secret is necessary if you are using the Authorization Code Flow.
   
   Client Assertion Signature Algorithm
   
   Signature algorithm to create JWT assertion as client authentication. In the case of JWT signed with private key or Client secret as jwt, it is required. If no algorithm is specified, the following algorithm is adapted. `RS256` is adapted in the case of JWT signed with private key. `HS256` is adapted in the case of Client secret as jwt.
   
   Client Assertion Audience
   
   The audience to use for the client assertion. The default value is the IDP’s token endpoint URL.
   
   Issuer
   
   Keycloak validates issuer claims, in responses from the IDP, against this value.
   
   Default Scopes
   
   A list of OIDC scopes Keycloak sends with the authentication request. The default value is `openid`. A space separates each scope.
   
   Prompt
   
   The prompt parameter in the OIDC specification. Through this parameter, you can force re-authentication and other options. See the specification for more details.
   
   Accepts prompt=none forward from client
   
   Specifies if the IDP accepts forwarded authentication requests containing the `prompt=none` query parameter. If a realm receives an auth request with `prompt=none`, the realm checks if the user is currently authenticated and returns a `login_required` error if the user has not logged in. When Keycloak determines a default IDP for the auth request (using the `kc_idp_hint` query parameter or having a default IDP for the realm), you can forward the auth request with `prompt=none` to the default IDP. The default IDP checks the authentication of the user there. Because not all IDPs support requests with `prompt=none`, Keycloak uses this switch to indicate that the default IDP supports the parameter before redirecting the authentication request.
   
   If the user is unauthenticated in the IDP, the client still receives a `login_required` error. If the user is authentic in the IDP, the client can still receive an `interaction_required` error if Keycloak must display authentication pages that require user interaction. This authentication includes required actions (for example, password change), consent screens, and screens set to display by the `first broker login` flow or `post broker login` flow.
   
   Requires short state parameter
   
   This switch needs to be enabled if identity provider does not support long value of the `state` parameter sent in the initial OIDC authentication request (EG. more than 100 characters). In this case, Keycloak will try to make shorter `state` parameter and may omit some client data to be sent in the initial request. This may result in the limited functionality in some very corner case scenarios (EG. in case that IDP redirects to Keycloak with the error in the OIDC authentication response, Keycloak might need to display error page instead of being able to redirect to the client in case that login session is expired).
   
   Validate Signatures
   
   Specifies if Keycloak verifies signatures on the external ID Token signed by this IDP. If **ON**, Keycloak must know the public key of the external OIDC IDP. For performance purposes, Keycloak caches the public key of the external OIDC identity provider.
   
   Use JWKS URL
   
   This switch is applicable if `Validate Signatures` is **ON**. If **Use JWKS URL** is **ON**, Keycloak downloads the IDP’s public keys from the JWKS URL. New keys download when the identity provider generates a new keypair. If **OFF**, Keycloak uses the public key (or certificate) from its database, so when the IDP keypair changes, import the new key to the Keycloak database as well.
   
   JWKS URL
   
   The URL pointing to the location of the IDP JWK keys. For more information, see the [JWK specification](https://datatracker.ietf.org/doc/html/rfc7517). If you use an external Keycloak as an IDP, you can use a URL such as [http://broker-keycloak:8180/realms/test/protocol/openid-connect/certs](http://broker-keycloak:8180/realms/test/protocol/openid-connect/certs) if your brokered Keycloak is running on [http://broker-keycloak:8180](http://broker-keycloak:8180) and its realm is `test`.
   
   Validating Public Key
   
   The public key in PEM format that Keycloak uses to verify external IDP signatures. This key applies if `Use JWKS URL` is **OFF**.
   
   Validating Public Key Id
   
   This setting applies if **Use JWKS URL** is **OFF**. This setting specifies the ID of the public key in PEM format. Because there is no standard way for computing key ID from the key, external identity providers can use different algorithms from what Keycloak uses. If this field’s value is not specified, Keycloak uses the validating public key for all requests, regardless of the key ID sent by the external IDP. When **ON**, this field’s value is the key ID used by Keycloak for validating signatures from providers and must match the key ID specified by the IDP.
   
   Forwarded query parameters
   
   Define the query parameters to be forwarded to an external AS from the initial authorization request sent to the authorization endpoint. Multiple parameters can be entered, separated by comma (,). The parameters available to forward are any non OpenID Connect/OAuth standard parameter or standard parameters that are available as a client note from the authentication session.
   
   Supports client assertions
   
   This setting enables support for using client assertions issued by the provider to authenticate clients. This requires to set Issuer and keys of this Identity provider. Keys can be set by setup of 'JWKS URL' or 'Validating public key' option.
   
   Allows client assertions to be re-used
   
   By default, a client assertion can not be used multiple times. If the client is not able to retrieve a new client assertion for each request this option can be enabled to allow re-use of the same client assertion.
   
   Allows Client ID as audience for assertions
   
   If enabled, the Client ID configured in the Identity Provider is the only valid audience for assertions used in Federated client authentication and in JWT Authorization Grants (Client Assertions and JWT Authorization Grant). The client ID is used instead of the token-url/issuer-url defined in the respective specifications. Note this behavior is not covered by any standard.

You can import all this configuration data by providing a URL or file that points to OpenID Provider Metadata. If you connect to a Keycloak external IDP, you can import the IDP settings from `<root>/realms/{realm-name}/.well-known/openid-configuration`. This link is a JSON document describing metadata about the IDP.

If you want to use [Json Web Encryption (JWE)](https://datatracker.ietf.org/doc/html/rfc7516) ID Tokens or UserInfo responses in the provider, the IDP needs to know the public key to use with Keycloak. The provider uses the [realm keys](#realm_keys) defined for the different encryption algorithms to decrypt the tokens. Keycloak provides a standard [JWKS endpoint](#con-server-oidc-uri-endpoints_server_administration_guide) which the IDP can use for downloading the keys automatically.

### [](#_identity_broker_oauth)OAuth v2 identity providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/oauth2.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Foauth2.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Foauth2.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak brokers identity providers based on the OAuth v2 protocol. These identity providers (IDPs) must support the [Authorization Code Flow](#con-oidc-auth-flows_server_administration_guide) defined in the specification to authenticate users and authorize access.

Procedure

1. Click **Identity Providers** in the menu.
2. From the `Add provider` list, select `OAuth v2`.
3. Enter your initial configuration options. See [General IDP Configuration](#_general-idp-config) for more information about configuration options.
   
   Table 3. OAuth2 settings   Configuration Description
   
   Authorization URL
   
   The authorization URL endpoint.
   
   Token URL
   
   The token URL endpoint.
   
   User Info URL
   
   An endpoint from where information about the user will be fetched from. When invoking this endpoint, Keycloak will send the request with the access token issued by the identity provider as a bearer token. As a result, it expects the response to be a JSON document with the claims that should be used to obtain user profile information like ID, username, email, and first and last names.
   
   Client Authentication
   
   Defines the Client Authentication method Keycloak uses with the Authorization Code Flow. In the case of JWT signed with a private key, Keycloak uses the realm private key. In the other cases, define a client secret. See the [Client Authentication specifications](https://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication) for more information.
   
   Client ID
   
   A realm acting as an OIDC client to the external IDP. The realm must have an OIDC client ID if you use the Authorization Code Flow to interact with the external IDP.
   
   Client Secret
   
   Client secret from an external [vault](#_vault-administration). This secret is necessary if you are using the Authorization Code Flow.
   
   Client Assertion Signature Algorithm
   
   Signature algorithm to create JWT assertion as client authentication. In the case of JWT signed with private key or Client secret as jwt, it is required. If no algorithm is specified, the following algorithm is adapted. `RS256` is adapted in the case of JWT signed with private key. `HS256` is adapted in the case of Client secret as jwt.
   
   Client Assertion Audience
   
   The audience to use for the client assertion. The default value is the IDP’s token endpoint URL.
   
   Default Scopes
   
   A space separated list of scopes Keycloak sends with the authentication request.
   
   Prompt
   
   The prompt parameter in the OIDC specification. Through this parameter, you can force re-authentication and other options. See the specification for more details.
   
   Accepts prompt=none forward from client
   
   Specifies if the IDP accepts forwarded authentication requests containing the `prompt=none` query parameter. If a realm receives an auth request with `prompt=none`, the realm checks if the user is currently authenticated and returns a `login_required` error if the user has not logged in. When Keycloak determines a default IDP for the auth request (using the `kc_idp_hint` query parameter or having a default IDP for the realm), you can forward the auth request with `prompt=none` to the default IDP. The default IDP checks the authentication of the user there. Because not all IDPs support requests with `prompt=none`, Keycloak uses this switch to indicate that the default IDP supports the parameter before redirecting the authentication request.
   
   If the user is unauthenticated in the IDP, the client still receives a `login_required` error. If the user is authentic in the IDP, the client can still receive an `interaction_required` error if Keycloak must display authentication pages that require user interaction. This authentication includes required actions (for example, password change), consent screens, and screens set to display by the `first broker login` flow or `post broker login` flow.
   
   Requires short state parameter
   
   This switch needs to be enabled if identity provider does not support long value of the `state` parameter sent in the initial OAuth2 authorization request (EG. more than 100 characters). In this case, Keycloak will try to make shorter `state` parameter and may omit some client data to be sent in the initial request. This may result in the limited functionality in some very corner case scenarios (EG. in case that IDP redirects to Keycloak with the error in the OAuth2 authorization response, Keycloak might need to display error page instead of being able to redirect to the client in case that login session is expired).

After the user authenticates to the identity provider and is redirected back to Keycloak, the broker will fetch the user profile information from the endpoint defined in the `User Info URL` setting. For that, Keycloak will invoke that endpoint using the access token issued by the identity provider as a bearer token. Even though the OAuth2 standard supports access tokens using a JWT format, this broker assumes access tokens are opaque and that user profile information should be obtained from a separate endpoint.

In order to map the claims from the JSON document returned by the user profile endpoint, you might want to set the following settings so that they are mapped to user attributes when federating the user:

Table 4. User profile claims   Configuration Description

ID Claim

The name of the claim from the JSON document returned by the user profile endpoint representing the user’s unique identifier. If not provided, defaults to `sub`.

Username Claim

The name of the claim from the JSON document returned by the user profile endpoint representing the user’s username. If not provided, defaults to `preferred_username`.

Email Claim

The name of the claim from the JSON document returned by the user profile endpoint representing the user’s email. If not provided, defaults to `email`.

Name Claim

The name of the claim from the JSON document returned by the user profile endpoint representing the user’s full name. If not provided, defaults to `name`.

Given name Claim

The name of the claim from the JSON document returned by the user profile endpoint representing the user’s given name. If not provided, defaults to `given_name`.

Family name Claim

The name of the claim from the JSON document returned by the user profile endpoint representing the user’s family name. If not provided, defaults to `family_name`.

You can import all this configuration data by providing a URL or file that points to the Authorization Server Metadata. If you connect to a Keycloak external IDP, you can import the IDP settings from `<root>/realms/{realm-name}/.well-known/openid-configuration`. This link is a JSON document describing metadata about the IDP.

### [](#saml-v2-0-identity-providers)SAML v2.0 Identity Providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/saml.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsaml.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsaml.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak can broker identity providers based on the SAML v2.0 protocol.

Procedure

1. Click **Identity Providers** in the menu.
2. From the `Add provider` list, select `SAML v2.0`.
   
   Add identity provider
   
   ![Add Identity Provider](./images/saml-add-identity-provider.png)
3. Enter your initial configuration options. See [General IDP Configuration](#_general-idp-config) for more information about configuration options.

Table 5. SAML Config   Configuration Description

Service Provider Entity ID

The SAML Entity ID that the remote Identity Provider uses to identify requests from this Service Provider. By default, this setting is set to the realms base URL `<root>/realms/{realm-name}`.

Identity Provider Entity ID

The Entity ID used to validate the Issuer for received SAML assertions. If empty, no Issuer validation is performed.

Single Sign-On Service URL

The SAML endpoint that starts the authentication process. If your SAML IDP publishes an IDP entity descriptor, the value of this field is specified there.

Artifact service URL

The SAML artifact resolution endpoint. If your SAML IDP publishes an IDP entity descriptor, the value of this field is specified there.

Single Logout Service URL

The SAML logout endpoint. If your SAML IDP publishes an IDP entity descriptor, the value of this field is specified there.

Backchannel Logout

Toggle this switch to **ON** if your SAML IDP supports back channel logout.

NameID Policy Format

The URI reference corresponding to a name identifier format. By default, Keycloak sets it to `urn:oasis:names:tc:SAML:2.0:nameid-format:persistent`.

Principal Type

Specifies which part of the SAML assertion will be used to identify and track external user identities. Can be either Subject NameID or SAML attribute (either by name or by friendly name). Subject NameID value can not be set together with 'urn:oasis:names:tc:SAML:2.0:nameid-format:transient' NameID Policy Format value.

Principal Attribute

If a Principal type is non-blank, this field specifies the name ("Attribute \[Name]") or the friendly name ("Attribute \[Friendly Name]") of the identifying attribute.

Allow create

Allow the external identity provider to create a new identifier to represent the principal.

HTTP-POST Binding Response

Controls the SAML binding in response to any SAML requests sent by an external IDP. When **OFF**, Keycloak uses Redirect Binding.

ARTIFACT Binding Response

Controls the SAML binding in response to any SAML requests sent by an external IDP. When **OFF**, Keycloak evaluates the HTTP-POST Binding Response configuration.

HTTP-POST Binding for AuthnRequest

Controls the SAML binding when requesting authentication from an external IDP. When **OFF**, Keycloak uses Redirect Binding.

Want AuthnRequests Signed

When **ON**, Keycloak uses the realm’s keypair to sign requests sent to the external SAML IDP.

Want Assertions Signed

Indicates whether this service provider expects a signed Assertion.

Want Assertions Encrypted

Indicates whether this service provider expects an encrypted Assertion.

Signature Algorithm

If **Want AuthnRequests Signed** is **ON**, the signature algorithm to use. Note that `SHA1` based algorithms are deprecated and may be removed in a future release. We recommend to use some more secure algorithm instead of `*_SHA1`. Also, with `*_SHA1` algorithms, verifying signatures do not work if the SAML identity provider (for example another instance of Keycloak) runs on Java 17 or higher.

Encryption Algorithm

Encryption algorithm, which is used by SAML IDP for encryption of SAML documents, assertions, or IDs. The corresponding decryption key for decrypt SAML document parts will be chosen based on this configured algorithm and should be available in realm keys for the encryption (ENC) usage. If the algorithm is not configured, any supported algorithm is allowed and a decryption key will be chosen based on the algorithm specified in SAML document itself.

SAML Signature Key Name

Signed SAML documents sent using POST binding contain the identification of signing key in `KeyName` element, which, by default, contains the Keycloak key ID. External SAML IDPs can expect a different key name. This switch controls whether `KeyName` contains: * `KEY_ID` - Key ID. * `CERT_SUBJECT` - the subject from the certificate corresponding to the realm key. Microsoft Active Directory Federation Services expect `CERT_SUBJECT`. * `NONE` - Keycloak omits the key name hint from the SAML message.

Force Authentication

The user must enter their credentials at the external IDP even when the user is already logged in.

Validate Signature

When **ON**, the realm expects SAML requests and responses from the external IDP to be digitally signed.

Metadata descriptor URL

External URL where Identity Provider publishes the `IDPSSODescriptor` metadata. This URL is used to download the Identity Provider certificates when the `Reload keys` or `Import keys` actions are clicked.

Use metadata descriptor URL

When **ON**, the certificates to validate signatures are automatically downloaded from the `Metadata descriptor URL` and cached in Keycloak. The SAML provider can validate signatures in two different ways. If a specific certificate is requested (usually in `POST` binding) and it is not in the cache, certificates are automatically refreshed from the URL. If all certificates are requested to validate the signature (`REDIRECT` binding) the refresh is only done after a max cache time. This maximum time can be specified in the descriptor itself, `cacheDuration` or `validUntil` attributes, or the cache provider defines one. See [public-key-storage](https://www.keycloak.org/server/all-provider-config) spi in the all provider config guide for more information about how the cache works.

When the option is **OFF**, the certificates in `Validating X509 Certificates` are used to validate signatures.

Validating X509 Certificates

The public certificates Keycloak uses to validate the signatures of SAML requests and responses from the external IDP when `Use metadata descriptor URL` is **OFF**. Multiple certificates can be entered separated by comma (`,`). The certificates can be re-imported from the `Metadata descriptor URL` clicking the `Import Keys` action in the identity provider page. The action downloads the current certificates in the metadata endpoint and assigns them to the config in this same option. You need to click `Save` to definitely store the re-imported certificates.

Sign Service Provider Metadata

When **ON**, Keycloak uses the realm’s key pair to sign the [SAML Service Provider Metadata descriptor](#_identity_broker_saml_sp_descriptor).

Pass subject

Controls if Keycloak forwards a `login_hint` query parameter to the IDP. Keycloak adds this field’s value to the login\_hint parameter in the AuthnRequest’s Subject so destination providers can pre-fill their login form.

Attribute Consuming Service Index

Identifies the attribute set to request to the remote IDP. Keycloak automatically adds the attributes mapped in the identity provider configuration to the autogenerated SP metadata document.

Attribute Consuming Service Name

A descriptive name for the set of attributes that are advertised in the autogenerated SP metadata document.

You can import all configuration data by providing a URL or a file pointing to the SAML IDP entity descriptor of the external IDP. If you are connecting to a Keycloak external IDP, you can import the IDP settings from the URL `<root>/realms/{realm-name}/protocol/saml/descriptor`. This link is an XML document describing metadata about the IDP. You can also import all this configuration data by providing a URL or XML file pointing to the external SAML IDP’s entity descriptor to connect to.

#### [](#_identity_broker_saml_requested_authncontext)Requesting specific AuthnContexts

Identity Providers facilitate clients specifying constraints on the authentication method verifying the user identity. For example, asking for MFA, Kerberos authentication, or security requirements. These constraints use particular AuthnContext criteria. A client can ask for one or more criteria and specify how the Identity Provider must match the requested AuthnContext, exactly, or by satisfying other equivalents.

You can list the criteria your Service Provider requires by adding ClassRefs or DeclRefs in the Requested AuthnContext Constraints section. Usually, you need to provide either ClassRefs or DeclRefs, so check with your Identity Provider documentation which values are supported. If no ClassRefs or DeclRefs are present, the Identity Provider does not enforce additional constraints.

Table 6. Requested AuthnContext Constraints   Configuration Description

Comparison

The method the Identity Provider uses to evaluate the context requirements. The available values are `Exact`, `Minimum`, `Maximum`, or `Better`. The default value is `Exact`.

AuthnContext ClassRefs

The AuthnContext ClassRefs describing the required criteria.

AuthnContext DeclRefs

The AuthnContext DeclRefs describing the required criteria.

#### [](#_identity_broker_saml_sp_descriptor)SP Descriptor

When you access the provider’s SAML SP metadata, look for the `Endpoints` item in the identity provider configuration settings. It contains a `SAML 2.0 Service Provider Metadata` link which generates the SAML entity descriptor for the Service Provider. You can download the descriptor or copy its URL and then import it into the remote Identity Provider.

This metadata is also available publicly by going to the following URL:

```
http[s]://{host:port}/realms/{realm-name}/broker/{broker-alias}/endpoint/descriptor
```

Ensure you save any configuration changes before accessing the descriptor.

#### [](#_identity_broker_saml_login_hint)Send subject in SAML requests

By default, a social button pointing to a SAML Identity Provider redirects the user to the following login URL:

```
http[s]://{host:port}/realms/${realm-name}/broker/{broker-alias}/login
```

Adding a query parameter named `login_hint` to this URL adds the parameter’s value to SAML request as a Subject attribute. If this query parameter is empty, Keycloak does not add a subject to the request.

Enable the "Pass subject" option to send the subject in SAML requests.

### [](#_identity_broker_spiffe)SPIFFE identity providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/spiffe.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fspiffe.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fspiffe.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

SPIFFE is **Preview** and is not fully supported. This feature is disabled by default.

To enable start the server with `--features=preview` or `--features=spiffe`

A SPIFFE identity provider supports authenticating clients with SPIFFE JWT SVIDs.

Each client must have a unique subject identifier (SPIFFE ID) for each realm.

SPIFFE JWT SVIDs can be re-used multiple times. As a security best practice, reduce the expiration time and make sure tokens can not be intercepted in communication between the client and Keycloak.

Procedure

1. Click **Identity Providers** in the menu.
2. From the `Add provider` list, select `SPIFFE`.
   
   Add SPIFFE provider
   
   ![Add SPIFFE Provider](./images/spiffe-add-identity-provider.png)
3. Enter your initial configuration options.
   
   Table 7. SPIFFE settings   Configuration Description
   
   Alias
   
   The alias for the identity provider is used to link a client to the provider
   
   SPIFFE Trust Domain
   
   The SPIFFE Trust domain (for example `spiffe://my-trust-domain`)
   
   SPIFFE Bundle Endpoint
   
   `https` URL for the SPIFFE Bundle Endpoint or the OpenID Connect JWKS endpoint where the SPIFFE servers public keys are exposed

### [](#_identity_broker_kubernetes)Kubernetes identity providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/kubernetes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fkubernetes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fkubernetes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

A Kubernetes identity provider supports authenticating clients with Kubernetes service account tokens.

This feature can be enabled and disabled with the feature `kubernetes-service-accounts` which is enabled by default. It depends on the feature `client-auth-federated` which is enabled by default as well.

The default issuer URL for a Kubernetes cluster is `https://kubernetes.default.svc.cluster.local`. You can discover this value by decoding a service account token to retrieve the value of the `iss` claim.

Keycloak must be able to invoke the endpoint `<ISSUER>/.well-known/openid-configuration` and additionally the JWKS endpoint returned in the well-known configuration. By default, these endpoints require authentication with a service account token. Keycloak will automatically use the token from `/var/run/secrets/kubernetes.io/serviceaccount/token` if available and the token issuer matches the configured issuer.

Each identity provider must have a unique issuer. Each client must also have a unique subject identifier for each issuer. As the subject identifier is built from the namespace and service account name, each client must have its own service account if multiple clients share a namespace.

As a security best practice, do not use the `default` service account in a namespace, as it is shared with all Pods in a namespace. Instead, create an individual service account for each client.

Kubernetes service account tokens can be reused multiple times. As a security best practice, reduce the expiration time and make sure tokens can not be intercepted in communication between the client and Keycloak.

Procedure

1. Click **Identity Providers** in the menu.
2. From the `Add provider` list, select `Kubernetes`.
3. Enter your initial configuration options or proceed with the defaults.
   
   Table 8. Kubernetes settings   Configuration Description
   
   Alias
   
   The alias for the identity provider is used to link a client to the provider
   
   Kubernetes Issuer URL
   
   The issuer URL of service account tokens. The URL `<ISSUER>.well-known/openid-configuration` must be available to Keycloak)
4. When you create a new realm with the preview feature `client-auth-federated` enabled, the client authentication flow is already configured correctly. For existing realms, add to the client authentication flow the execution of **Signed JWT - Federated** as an alternative step. As built-in flows can not be updated, and if the default flow is your default, you will first need to duplicate the existing clients
5. For each confidential OpenID Connect client that should authenticate via this provider:
   
   1. Change in the **Credentials** tab the **Client Authenticator** to **Signed JWT - Federated**.
   2. As **Identity provider**, enter the alias of the Kubernetes identity provider added in step 3.
   3. As **Federated Subject**, enter the subject identifier as issued by Kubernetes. This is usually `system:serviceaccount:<namespace>:<serviceaccount>`.
6. For the Pod in Kubernetes add a service account:
   
   ```
   apiVersion: v1
   kind: Pod
   ...
   spec:
     serviceAccountName: <serviceaccount>
     ...
     volumes:
     - name: aud-token
       projected:
         defaultMode: 420
         sources:
         - serviceAccountToken:
             audience: https://example.com:8443/realms/test (1)
             expirationSeconds: 600 (2)
             path: my-aud-token
   ```
   
   1. Issuer URL of the Keycloak realm.
   2. Maximum time allowed by Kubernetes and Keycloak is 3600 seconds
   
   To verify your setup, assuming the client has a service account configured:
7. In the Pod, retrieve the token
   
   ```
   cat /var/run/secrets/tokens/my-aud-token
   ```
8. Use the token as the client credentials
   
   ```
   curl -k https://example.com:8443/realms/<realm>/protocol/openid-connect/token \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     --data-urlencode grant_type=client_credentials \
     --data-urlencode client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer \
     --data-urlencode client_assertion=<token>
   ```
9. Verify the response
   
   ```
   {
     "access_token": "ey...bw",
     "expires_in": 300,
     "refresh_expires_in": 0,
     "token_type": "Bearer",
     "not-before-policy": 0,
     "scope": "profile email"
   }
   ```

While the service account functionality is helpful to test that the setup is working, disabled this feature after the test if it is not needed for the client to follow the least-privileges security best practice.

### [](#_client_suggested_idp)Client-suggested Identity Provider

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/suggested.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsuggested.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsuggested.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

OIDC applications can bypass the Keycloak login page by hinting at the identity provider they want to use. You can enable this by setting the `kc_idp_hint` query parameter in the Authorization Code Flow authorization endpoint.

With Keycloak OIDC client adapters, you can specify this query parameter when you access a secured resource in the application.

For example:

```
GET /myapplication.com?kc_idp_hint=facebook HTTP/1.1
Host: localhost:8080
```

In this case, your realm must have an identity provider with a `facebook` alias. If this provider does not exist, the login form is displayed.

If you are using the JavaScript adapter, you can also achieve the same behavior as follows:

```
const keycloak = new Keycloak({
        url: "http://keycloak-server",
        realm: "my-realm",
        clientId: "my-app"
);

await keycloak.createLoginUrl({
        idpHint: 'facebook'
});
```

With the `kc_idp_hint` query parameter, the client can override the default identity provider if you configure one for the `Identity Provider Redirector` authenticator. The client can disable the automatic redirecting by setting the `kc_idp_hint` query parameter to an empty value.

### [](#_mappers)Mapping claims and assertions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/mappers.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fmappers.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fmappers.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can import the SAML and OpenID Connect metadata, provided by the external IDP you are authenticating with, into the realm. After importing, you can extract user profile metadata and other information, so you can make it available to your applications.

Each user logging into your realm using an external identity provider has an entry in the local Keycloak database, based on the metadata from the SAML or OIDC assertions and claims.

Procedure

1. Click **Identity Providers** in the menu.
2. Select one of the identity providers in the list.
3. Click the **Mappers** tab.
   
   Identity provider mappers
   
   ![identity provider mappers](./images/identity-provider-mappers.png)
4. Click **Add mapper**.
   
   Identity provider mapper
   
   ![identity provider mapper](./images/identity-provider-mapper.png)
5. Select a value for **Sync Mode Override**. The mapper updates user information when users log in repeatedly according to this setting.
   
   1. Select **legacy** to use the behavior of the previous Keycloak version.
   2. Select **import** to import data from when the user was first created in Keycloak during the first login to Keycloak with a particular identity provider.
   3. Select **force** to update user data at each user login.
   4. Select **inherit** to use the sync mode configured in the identity provider. All other options will override this sync mode.
6. Select a mapper from the **Mapper Type** list. Hover over the **Mapper Type** for a description of the mapper and configuration to enter for the mapper.
7. Click **Save**.

For JSON-based claims, you can use dot notation for nesting and square brackets to access array fields by index. For example, `contact.address[0].country`.

To investigate the structure of user profile JSON data provided by social providers, you can enable the `DEBUG` level logger `org.keycloak.social.user_profile_dump` when starting the server.

### [](#available-user-session-data)Available user session data

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/session-data.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsession-data.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fsession-data.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

After a user login from an external IDP, Keycloak stores user session note data that you can access. This data can be propagated to the client requesting log in using the token or SAML assertion passed back to the client using an appropriate client mapper.

identity\_provider

The IDP alias of the broker used to perform the login.

identity\_provider\_identity

The IDP username of the currently authenticated user. Often, but not always, the same as the Keycloak username. For example, Keycloak can link a user john\` to a Facebook user `john123@gmail.com`. In that case, the value of the user session note is `john123@gmail.com`.

You can use a [Protocol Mapper](#_protocol-mappers) of type `User Session Note` to propagate this information to your clients.

### [](#_identity_broker_first_login)First login flow

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/first-login-flow.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Ffirst-login-flow.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Ffirst-login-flow.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When users log in through identity brokering, Keycloak imports and links aspects of the user within the realm’s local database. When Keycloak successfully authenticates users through an external identity provider, two situations can exist:

- Keycloak has already imported and linked a user account with the authenticated identity provider account. In this case, Keycloak authenticates as the existing user and redirects back to the application.
- No account exists for this user in Keycloak. Usually, you register and import a new account into the Keycloak database, but there may be an existing Keycloak account with the same email address. Automatically linking the existing local account to the external identity provider is a potential security hole. You cannot always trust the information you get from the external identity provider.

Different organizations have different requirements when dealing with some of these situations. With Keycloak, you can use the `First Login Flow` option in the IDP settings to choose a [workflow](#_authentication-flows) for a user logging in from an external IDP for the first time. By default, the `First Login Flow` option points to the `first broker login` flow, but you can use your flow or different flows for different identity providers.

The flow is in the Admin Console under the **Authentication** tab. When you choose the `First Broker Login` flow, you see the authenticators used by default. You can re-configure the existing flow. For example, you can disable some authenticators, mark some of them as `required`, or configure some authenticators.

You can also create a new authentication flow, write your own Authenticator implementations, and use it in your flow. See [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more information.

#### [](#default-first-login-flow-authenticators)Default first login flow authenticators

Review Profile

- This authenticator displays the profile information page, so the users can review their profile that Keycloak retrieves from an identity provider.
- You can set the `Update Profile On First Login` option in the **Actions** menu.
- When **ON**, users are presented with the profile page requesting additional information to federate the user’s identities.
- When **missing**, users are presented with the profile page if the identity provider does not provide mandatory information, such as email, first name, or last name.
- When **OFF**, the profile page does not display unless the user clicks in a later phase on the `Review profile info` link in the page displayed by the `Confirm Link Existing Account` authenticator.

Create User If Unique

This authenticator checks if there is already an existing Keycloak account with the same email or username like the account from the identity provider. If it’s not, then the authenticator just creates a new local Keycloak account and links it with the identity provider and the whole flow is finished. Otherwise it goes to the next `Handle Existing Account` subflow. If you always want to ensure that there is no duplicated account, you can mark this authenticator as `REQUIRED`. In this case, the user will see the error page if there is an existing Keycloak account and the user will need to link the identity provider account through Account management.

- This authenticator verifies that there is already a Keycloak account with the same email or username as the identity provider’s account.
- If an account does not exist, the authenticator creates a local Keycloak account, links this account with the identity provider, and terminates the flow.
- If an account exists, the authenticator implements the next `Handle Existing Account` sub-flow.
- To ensure there is no duplicated account, you can mark this authenticator as `REQUIRED`. The user sees the error page if a Keycloak account exists, and users must link their identity provider account through Account management.

Confirm Link Existing Account

- On the information page, users see a Keycloak account with the same email. Users can review their profile again and use a different email or username. The flow restarts and goes back to the `Review Profile` authenticator.
- Alternatively, users can confirm that they want to link their identity provider account with their existing Keycloak account.
- Disable this authenticator if you do not want users to see this confirmation page and go straight to linking identity provider account by email verification or re-authentication.

Verify Existing Account By Email

- This authenticator is `ALTERNATIVE` by default. Keycloak uses this authenticator if the realm has an SMTP setup configured.
- The authenticator sends an email to users to confirm that they want to link the identity provider with their Keycloak account.
- Disable this authenticator if you do not want to confirm linking by email, but want users to reauthenticate with their password.

Verify Existing Account By Re-authentication

- Use this authenticator if the email authenticator is not available. For example, you have not configured SMTP for your realm. This authenticator displays a login screen for users to authenticate to link their Keycloak account with the Identity Provider.
- Users can also re-authenticate with another identity provider already linked to their Keycloak account.
- You can force users to use OTP. Otherwise, it is optional and used if you have set OTP for the user account.

#### [](#automatically-link-existing-first-login-flow)Automatically link existing first login flow

The AutoLink authenticator is dangerous in a generic environment where users can register themselves using arbitrary usernames or email addresses. Do not use this authenticator unless you are carefully curating user registration and assigning usernames and email addresses.

To configure a first login flow that links users automatically without prompting, create a new flow with the following two authenticators:

Create User If Unique

This authenticator ensures Keycloak handles unique users. Set the authenticator requirement to **Alternative**.

Automatically Set Existing User

This authenticator sets an existing user to the authentication context without verification. Set the authenticator requirement to "Alternative".

This setup is the simplest setup available, but it is possible to use other authenticators. For example:

- You can add the Review Profile authenticator to the beginning of the flow if you want end users to confirm their profile information.
- You can add authentication mechanisms to this flow, forcing a user to verify their credentials. Adding authentication mechanisms requires a complex flow. For example, you can set the "Automatically Set Existing User" and "Password Form" as "Required" in an "Alternative" sub-flow.

#### [](#_disabling_automatic_user_creation)Disabling automatic user creation

The Default first login flow looks up the Keycloak account matching the external identity and offers to link them. If no matching Keycloak account exists, the flow automatically creates one.

This default behavior may be unsuitable for some setups. One example is when you use a read-only LDAP user store, where all users are pre-created. In this case, you must switch off automatic user creation.

To disable user creation:

Procedure

1. Click **Authentication** in the menu.
2. Select **First Broker Login** from the list.
3. Set **Create User If Unique** to **DISABLED**.
4. Set **Confirm Link Existing Account** to **DISABLED**.

This configuration also implies that Keycloak itself won’t be able to determine which internal account would correspond to the external identity. Therefore, the `Verify Existing Account By Re-authentication` authenticator will ask the user to provide both username and password.

Enabling or disabling user creation by identity provider is completely independent on the realm [User Registration switch](#con-user-registration_server_administration_guide). You can have enabled user-creation by identity provider and at the same time disabled user self-registration in the realm login settings or vice-versa.

#### [](#_detect_existing_user_first_login_flow)Detect existing user first login flow

In order to configure a first login flow in which:

- only users already registered in this realm can log in,
- users are automatically linked without being prompted,

create a new flow with the following two authenticators:

Detect Existing Broker User

This authenticator ensures that unique users are handled. Set the authenticator requirement to `REQUIRED`.

Automatically Set Existing User

Automatically sets an existing user to the authentication context without any verification. Set the authenticator requirement to `REQUIRED`.

You have to set the `First Login Flow` of the identity provider configuration to that flow. You could set the also set `Sync Mode` to `force` if you want to update the user profile (Last Name, First Name…​) with the identity provider attributes.

This flow can be used if you want to delegate the identity to other identity providers (such as GitHub, Facebook …​) but you want to manage which users that can log in.

With this configuration, Keycloak is unable to determine which internal account corresponds to the external identity. The **Verify Existing Account By Re-authentication** authenticator asks the provider for the username and password.

#### [](#_override_existing_broker_link)Override existing broker link

When an another account needs to be linked to the same Keycloak account within the same identity provider, you can configure the following authenticator.

Confirm Override Existing Link

This authenticator will detect the existing broker link for the user and display a confirmation page to confirm overriding the existing broker link. Set the authenticator requirement to REQUIRED.

A typical use of this authenticator is a scenario such as the following:

- For example, consider a Keycloak user `john` with the email `john@gmail.com`. That user is linked to the identity provider `google` with the `google` username `john@gmail.com` .
- Then for instance Keycloak user `john` creates new Google account with email `john-new@gmail.com`
- Then during login to Keycloak, the user authenticated to the identity provider `google` with a new username such as `john-new@gmail.com`, which is not linked to any Keycloak account yet (as Keycloak account `john` is still linked with the `google` user `john@gmail.com`) and hence the first-broker-login flow is triggered.
- During first-broker-login, the Keycloak user `john` is authenticated somehow (either by default first-broker-login re-authentication or for instance by authenticator like `Detect existing broker user`)
- Now with this authenticator in the authentication flow, it is possible to override the IDP link to the `google` identity provider of Keycloak user `john` with the new `google` link to `google` user `john-new@gmail.com` after user `john` confirms this.

When creating authentication flows with this authenticator, make sure to add this authenticator once other authenticators that are already established the Keycloak user by other means (either by re-authentication or after `Detect existing broker user` as mentioned above.

### [](#_identity_broker_post_login_flow)Post login flow

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/post-login-flow.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fpost-login-flow.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Fpost-login-flow.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Post login flow is useful for the situations when you want to trigger some additional authentication actions after every login with the particular identity provider. For example, you may want to trigger 2-factor authentication after every login of Keycloak to `Facebook` because `Facebook` does not provide 2-factor authentication during its login.

Once you setup the authentication flow with the needed steps, set it as `Post login flow` when configuring the identity provider.

#### [](#post-login-flow-examples)Post login flow examples

##### [](#requesting-2-factor-authentication-after-identity-provider-login)Requesting 2-factor authentication after identity provider login

The easiest way is to enforce authentication with one particular 2-factor method. For example, when requesting OTP, the flow can look like this with only a single authenticator configured. This type of flow asks the user to configure the OTP during the first login with the identity provide when the user does not have OTP set on the account.

2FA post login flow with OTP

![Post login OTP](./images/post-login-flow-otp.png)

The more complex setup can include multiple 2-factor authentication methods configured as `ALTERNATIVE`. In this case, make sure that the user is requested to setup one of the methods if that user does not yet have any 2-factor authentication configured on the account. This could be done as follows:

- Make sure that one of the 2-factor methods is configured as `REQUIRED` in the [First login flow](#_identity_broker_first_login). This method can works if you expect all your users to be registered by the identity provider login.
- Wrap the 2-factor methods as `ALTERNATIVE` into a conditional subflow such as one called `2FA` and create another conditional subflow such as one called `OTP if no 2FA`, which will be triggered only if the previous subflow was not executed and will ask user to add one of the 2-factor methods (for example, OTP). The example of a similar flow configuration is provided in the [Conditions section of the Authentication flows chapter](#_conditional-2fa-otp-default).

#### [](#requesting-additional-authentication-steps-for-the-dedicated-clients)Requesting additional authentication steps for the dedicated clients

In some cases, a client or group of clients may need to perform some additional steps after identity provider login. The following is an example of a flow that prescribes that when the client scope `foo` is requested, the user is required to authenticate with the OTP after identity provider login.

2FA post login flow with client scope and OTP

![Post login with client scope and OTP](./images/post-login-flow-client-scope.png)

This is an example of configuring the `Condition - client scope` for requesting the specified client scope.

2FA post login flow client scope configuration

![Post login flow client scope configuration](./images/post-login-flow-client-scope-config.png)

The requested clients need to have this client scope set on them either as default or as optional. In the latter case, the flow is executed only if the client scope is requested by the client (for example, by the `scope` parameter in the case of OIDC/OAuth2 client logins).

### [](#retrieving-external-idp-tokens)Retrieving external IDP tokens

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/tokens.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Ftokens.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Ftokens.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

With Keycloak, you can store tokens and responses from the authentication process with the external IDP using the `Store Token` configuration option on the IDP’s settings page.

Application code can retrieve these tokens and responses to import extra user information or to request the external IDP securely. For example, an application can use the Google token to use other Google services and REST APIs. To retrieve a token for a particular identity provider, send a request as follows:

```
GET /realms/{realm-name}/broker/{provider_alias}/token HTTP/1.1
Host: localhost:8080
Authorization: Bearer <KEYCLOAK ACCESS TOKEN>
```

An application must authenticate with Keycloak and receive an access token. This access token must have the `broker` client-level role `read-token` set, so the user must have a role mapping for this role, and the client application must have that role within its scope. In this case, since you are accessing a protected service in Keycloak, send the access token issued by Keycloak during the user authentication. You can assign this role to newly imported users in the broker configuration page by setting the **Stored Tokens Readable** switch to **ON**.

These external tokens can be re-established by logging in again through the provider or using the client-initiated account linking API.

### [](#identity-broker-logout)Identity broker logout

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/identity-broker/logout.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fidentity-broker%2Flogout.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fidentity-broker%2Flogout.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When logging out, Keycloak sends a request to the external identity provider that is used to log in initially and logs the user out of this identity provider.

## [](#sso-protocols)SSO protocols

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sso-protocols.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsso-protocols.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsso-protocols.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

This section discusses authentication protocols, the Keycloak authentication server and how applications, secured by the Keycloak authentication server, interact with these protocols.

### [](#con-oidc_server_administration_guide)OpenID Connect

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sso-protocols/con-oidc.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-oidc.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-oidc.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

[OpenID Connect](https://openid.net/developers/how-connect-works/) (OIDC) is an authentication protocol that is an extension of [OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749).

OAuth 2.0 is a framework for building authorization protocols and is incomplete. OIDC, however, is a full authentication and authorization protocol that uses the [Json Web Token](https://www.jwt.io/) (JWT) standards. The JWT standards define an identity token JSON format and methods to digitally sign and encrypt data in a compact and web-friendly way.

In general, OIDC implements two use cases. The first case is an application requesting that a Keycloak server authenticates a user. Upon successful login, the application receives an *identity token* and an *access token*. The *identity token* contains user information including user name, email, and profile information. The realm digitally signs the *access token* which contains access information (such as user role mappings) that applications use to determine the resources users can access in the application.

The second use case is a client accessing remote services.

- The client requests an *access token* from Keycloak to invoke on remote services on behalf of the user.
- Keycloak authenticates the user and asks the user for consent to grant access to the requesting client.
- The client receives the *access token* which is digitally signed by the realm.
- The client makes REST requests on remote services using the *access token*.
- The remote REST service extracts the *access token*.
- The remote REST service verifies the tokens signature.
- The remote REST service decides, based on access information within the token, to process or reject the request.

#### [](#con-oidc-auth-flows_server_administration_guide)OIDC auth flows

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sso-protocols/con-oidc-auth-flows.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-oidc-auth-flows.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-oidc-auth-flows.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

OIDC has several methods, or flows, that clients or applications can use to authenticate users and receive *identity* and *access* tokens. The method depends on the type of application or client requesting access.

##### [](#_oidc-auth-flows-authorization)Authorization Code Flow

The Authorization Code Flow is a browser-based protocol and suits authenticating and authorizing browser-based applications. It uses browser redirects to obtain *identity* and *access* tokens.

1. A user connects to an application using a browser. The application detects the user is not logged into the application.
2. The application redirects the browser to Keycloak for authentication.
3. The application passes a callback URL as a query parameter in the browser redirect. Keycloak uses the parameter upon successful authentication.
4. Keycloak authenticates the user and creates a one-time, short-lived, temporary code.
5. Keycloak redirects to the application using the callback URL and adds the temporary code as a query parameter in the callback URL.
6. The application extracts the temporary code and makes a background REST invocation to Keycloak to exchange the code for an *identity* and *access* and *refresh* token. To prevent replay attacks, the temporary code cannot be used more than once.

A system is vulnerable to a stolen token for the lifetime of that token. For security and scalability reasons, access tokens are generally set to expire quickly so subsequent token requests fail. If a token expires, an application can obtain a new access token using the additional *refresh* token sent by the login protocol.

*Confidential* clients provide client secrets when they exchange the temporary codes for tokens. *Public* clients are not required to provide client secrets. *Public* clients are secure when HTTPS is strictly enforced and redirect URIs registered for the client are strictly controlled. HTML5/JavaScript clients have to be *public* clients because there is no way to securely transmit the client secret to HTML5/JavaScript clients. For more details, see the [Managing Clients](#assembly-managing-clients_server_administration_guide) chapter.

Keycloak also supports the [Proof Key for Code Exchange](https://datatracker.ietf.org/doc/html/rfc7636) specification.

##### [](#_oidc-auth-flows-implicit)Implicit Flow

The Implicit Flow is a browser-based protocol. It is similar to the Authorization Code Flow but with fewer requests and no refresh tokens.

The possibility exists of *access* tokens leaking in the browser history when tokens are transmitted via redirect URIs (see below).

Also, this flow does not provide clients with refresh tokens. Therefore, access tokens have to be long-lived or users have to re-authenticate when they expire.

We do not advise using this flow. This flow is supported because it is in the OIDC and OAuth 2.0 specification.

The protocol works as follows:

1. A user connects to an application using a browser. The application detects the user is not logged into the application.
2. The application redirects the browser to Keycloak for authentication.
3. The application passes a callback URL as a query parameter in the browser redirect. Keycloak uses the query parameter upon successful authentication.
4. Keycloak authenticates the user and creates an *identity* and *access* token. Keycloak redirects to the application using the callback URL and additionally adds the *identity* and *access* tokens as a query parameter in the callback URL.
5. The application extracts the *identity* and *access* tokens from the callback URL.

##### [](#_oidc-auth-flows-direct)Resource owner password credentials grant (Direct Access Grants)

*Direct Access Grants* are used by REST clients to obtain tokens on behalf of users. It is a HTTP POST request that contains:

- The credentials of the user. The credentials are sent within form parameters.
- The id of the client.
- The clients secret (if it is a confidential client).

The HTTP response contains the *identity*, *access*, and *refresh* tokens.

##### [](#_client_credentials_grant)Client credentials grant

The *Client Credentials Grant* creates a token based on the metadata and permissions of a service account associated with the client instead of obtaining a token that works on behalf of an external user. *Client Credentials Grants* are used by REST clients.

See the [Service Accounts](#_service_accounts) chapter for more information.

##### [](#_refresh_token_grant)Refresh token grant

By default, Keycloak returns refresh tokens in the token responses from most of the flows. Some exceptions are implicit flow or client credentials grant described above.

Refresh token is tied to the user session of the SSO browser session and can be valid for the lifetime of the user session. However, that client should send a refresh-token request at least once per specified interval. Otherwise, the session can be considered "idle" and can expire. See the [timeouts section](#_timeouts) for more information.

Keycloak supports [offline tokens](#_offline-access), which can be used typically when client needs to use refresh token even if corresponding browser SSO session is already expired.

###### [](#_refresh_token_rotation)Refresh token rotation

It is possible to specify that the refresh token is considered invalid once it is used. This means that client must always save the refresh token from the last refresh response because older refresh tokens, which were already used, would not be considered valid anymore by Keycloak. This is possible to set with the use of *Revoke Refresh token* option as specified in the [timeouts section](#_timeouts).

Keycloak also supports the situation that no refresh token rotation exists. In this case, a refresh token is returned during login, but subsequent responses from refresh-token requests will not return new refresh tokens. This practice is recommended for instance in the **FAPI 2 draft specification** and **FAPI 2 final specification** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section. In Keycloak, it is possible to skip refresh token rotation with the use of [client policies](#_client_policies). You can add executor `suppress-refresh-token-rotation` to some client profile and configure client policy to specify for which clients would be the profile triggered, which means that for those clients the refresh token rotation is going to be skipped.

##### [](#device-authorization-grant)Device authorization grant

This is used by clients running on internet-connected devices that have limited input capabilities or lack a suitable browser. Here’s a brief summary of the protocol:

1. The application requests Keycloak a device code and a user code. Keycloak creates a device code and a user code. Keycloak returns a response including the device code and the user code to the application.
2. The application provides the user with the user code and the verification URI. The user accesses a verification URI to be authenticated by using another browser. You could define a short verification\_uri that will be redirected to Keycloak verification URI (/realms/realm\_name/device) outside Keycloak - for example in a proxy.
3. The application repeatedly polls Keycloak to find out if the user completed the user authorization. If user authentication is complete, the application exchanges the device code for an *identity*, *access* and *refresh* token.

##### [](#_client_initiated_backchannel_authentication_grant)Client initiated backchannel authentication grant

This feature is used by clients who want to initiate the authentication flow by communicating with the OpenID Provider directly without redirect through the user’s browser like OAuth 2.0’s authorization code grant. Here’s a brief summary of the protocol:

1. The client requests Keycloak an auth\_req\_id that identifies the authentication request made by the client. Keycloak creates the auth\_req\_id.
2. After receiving this auth\_req\_id, this client repeatedly needs to poll Keycloak to obtain an Access Token, Refresh Token and ID Token from Keycloak in return for the auth\_req\_id until the user is authenticated.

An administrator can configure Client Initiated Backchannel Authentication (CIBA) related operations as `CIBA Policy` per realm.

Also please refer to other places of Keycloak documentation like **Backchannel Authentication Endpoint** and **Client Initiated Backchannel Authentication Grant** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section.

###### [](#ciba-policy)CIBA Policy

An administrator carries out the following operations on the `Admin Console` :

- Open the `Authentication → CIBA Policy` tab.
- Configure items and click `Save`.

The configurable items and their description follow.

  Configuration Description

Backchannel Token Delivery Mode

Specifying how the CD (Consumption Device) gets the authentication result and related tokens. There are three modes, "poll", "ping" and "push". Keycloak only supports "poll". The default setting is "poll". This configuration is required. For more details, see [CIBA Specification](https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html#rfc.section.5).

Expires In

The expiration time of the "auth\_req\_id" in seconds since the authentication request was received. The default setting is 120. This configuration is required. For more details, see [CIBA Specification](https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html#successful_authentication_request_acknowdlegment).

Interval

The interval in seconds the CD (Consumption Device) needs to wait for between polling requests to the token endpoint. The default setting is 5. This configuration is optional. For more details, see [CIBA Specification](https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html#successful_authentication_request_acknowdlegment).

Authentication Requested User Hint

The way of identifying the end-user for whom authentication is being requested. The default setting is "login\_hint". There are three modes, "login\_hint", "login\_hint\_token" and "id\_token\_hint". Keycloak only supports "login\_hint". This configuration is required. For more details, see [CIBA Specification](https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html#rfc.section.7.1).

###### [](#provider-setting)Provider Setting

The CIBA grant uses the following two providers.

1. Authentication Channel Provider: provides the communication between Keycloak and the entity that actually authenticates the user via AD (Authentication Device).
2. User Resolver Provider: get `UserModel` of Keycloak from the information provided by the client to identify the user.

Keycloak has both default providers. However, the administrator needs to set up Authentication Channel Provider like this:

```
kc.[sh|bat] start --spi-ciba-auth-channel--ciba-http-auth-channel--http-authentication-channel-uri=https://backend.internal.example.com
```

The configurable items and their description follow.

  Configuration Description

http-authentication-channel-uri

Specifying URI of the entity that actually authenticates the user via AD (Authentication Device).

###### [](#authentication-channel-provider)Authentication Channel Provider

CIBA standard document does not specify how to authenticate the user by AD. Therefore, it might be implemented at the discretion of products. Keycloak delegates this authentication to an external authentication entity. To communicate with the authentication entity, Keycloak provides Authentication Channel Provider.

Its implementation of Keycloak assumes that the authentication entity is under the control of the administrator of Keycloak so that Keycloak trusts the authentication entity. It is not recommended to use the authentication entity that the administrator of Keycloak cannot control.

Authentication Channel Provider is provided as SPI provider so that users of Keycloak can implement their own provider in order to meet their environment. Keycloak provides its default provider called HTTP Authentication Channel Provider that uses HTTP to communicate with the authentication entity.

If a user of Keycloak wants to use the HTTP Authentication Channel Provider, they need to know its contract between Keycloak and the authentication entity consisting of the following two parts.

Authentication Delegation Request/Response

Keycloak sends an authentication request to the authentication entity.

Authentication Result Notification/ACK

The authentication entity notifies the result of the authentication to Keycloak.

Authentication Delegation Request/Response consists of the following messaging.

Authentication Delegation Request

The request is sent from Keycloak to the authentication entity to ask it for user authentication by AD.

```
POST [delegation_reception]
```

- Headers

   Name Value Description

Content-Type

application/json

The message body is json formatted.

Authorization

Bearer \[token]

The \[token] is used when the authentication entity notifies the result of the authentication to Keycloak.

- Parameters

   Type Name Description

Path

delegation\_reception

The endpoint provided by the authentication entity to receive the delegation request

- Body

  Name Description

login\_hint

It tells the authentication entity who is authenticated by AD.  
By default, it is the user’s "username".  
This field is required and was defined by CIBA standard document.

scope

It tells which scopes the authentication entity gets consent from the authenticated user.  
This field is required and was defined by CIBA standard document.

is\_consent\_required

It shows whether the authentication entity needs to get consent from the authenticated user about the scope.  
This field is required.

binding\_message

Its value is intended to be shown in both CD and AD’s UI to make the user recognize that the authentication by AD is triggered by CD.  
This field is optional and was defined by CIBA standard document.

acr\_values

It tells the requesting Authentication Context Class Reference from CD.  
This field is optional and was defined by CIBA standard document.

Authentication Delegation Response

The response is returned from the authentication entity to Keycloak to notify that the authentication entity received the authentication request from Keycloak.

- Responses

  HTTP Status Code Description

201

It notifies Keycloak of receiving the authentication delegation request.

Authentication Result Notification/ACK consists of the following messaging.

Authentication Result Notification

The authentication entity sends the result of the authentication request to Keycloak.

```
POST /realms/[realm]/protocol/openid-connect/ext/ciba/auth/callback
```

- Headers

   Name Value Description

Content-Type

application/json

The message body is json formatted.

Authorization

Bearer \[token]

The \[token] must be the one the authentication entity has received from Keycloak in Authentication Delegation Request.

- Parameters

   Type Name Description

Path

realm

The realm name

- Body

  Name Description

status

It tells the result of user authentication by AD.  
It must be one of the following status.  
SUCCEED : The authentication by AD has been successfully completed.  
UNAUTHORIZED : The authentication by AD has not been completed.  
CANCELLED : The authentication by AD has been cancelled by the user.

Authentication Result ACK

The response is returned from Keycloak to the authentication entity to notify Keycloak received the result of user authentication by AD from the authentication entity.

- Responses

  HTTP Status Code Description

200

It notifies the authentication entity of receiving the notification of the authentication result.

###### [](#user-resolver-provider)User Resolver Provider

Even if the same user, its representation may differ in each CD, Keycloak and the authentication entity.

For CD, Keycloak and the authentication entity to recognize the same user, this User Resolver Provider converts their own user representations among them.

User Resolver Provider is provided as SPI provider so that users of Keycloak can implement their own provider in order to meet their environment. Keycloak provides its default provider called Default User Resolver Provider that has the following characteristics.

- Only support `login_hint` parameter and is used as default.
- `username` of UserModel in Keycloak is used to represent the user on CD, Keycloak and the authentication entity.

#### [](#_oidc-logout)OIDC Logout

OIDC has four specifications relevant to logout mechanisms:

1. [Session Management](https://openid.net/specs/openid-connect-session-1_0.html)
2. [RP-Initiated Logout](https://openid.net/specs/openid-connect-rpinitiated-1_0.html)
3. [Front-Channel Logout](https://openid.net/specs/openid-connect-frontchannel-1_0.html)
4. [Back-Channel Logout](https://openid.net/specs/openid-connect-backchannel-1_0.html)

Again since all of this is described in the OIDC specification we will only give a brief overview here.

##### [](#_oidc-logout-session-management)Session Management

This is a browser-based logout. The application obtains session status information from Keycloak at a regular basis. When the session is terminated at Keycloak the application will notice and trigger its own logout.

This is useful especially for the browser-based applications, such as when the application is secured by the [Keycloak Javascript adapter](https://www.keycloak.org/securing-apps/javascript-adapter) As a result, logout from one browser tab can trigger automatic logout from all other browser tabs with applications secured by the javascript adapter.

##### [](#rp-initiated-logout)RP-Initiated Logout

This is also a browser-based logout where the logout starts by redirecting the user to a specific endpoint at Keycloak. This redirect usually happens when the user clicks the `Log Out` link on the page of some application, which previously used Keycloak to authenticate the user.

Once the user is redirected to the logout endpoint, Keycloak is going to send logout requests to clients attached to the current browser SSO session to let them invalidate their local user sessions, and potentially redirect the user to some URL once the logout process is finished. The user might be optionally requested to confirm the logout in case the `id_token_hint` parameter was not used. If the client has [Logout confirmation](#_logout-confirmation) enabled, Keycloak renders a confirmation page after a successful logout informing the user that they are logged out. When a valid `post_logout_redirect_uri` is provided, this page includes an option to continue to that URL. After logout, the user is automatically redirected to the specified `post_logout_redirect_uri` as long as it is provided as a parameter. Note that you need to include either the `client_id` or `id_token_hint` parameter in case the `post_logout_redirect_uri` is included. Also the `post_logout_redirect_uri` parameter needs to match one of the `Valid Post Logout Redirect URIs` specified in the client configuration.

Depending on the client configuration, logout requests can be sent to clients through the front-channel or through the back-channel. For the frontend browser clients, which rely on the Session Management described in the previous section, Keycloak does not need to send any logout requests to them; these clients automatically detect that SSO session in the browser is logged out.

##### [](#front-channel-logout)Front-channel Logout

To configure clients to receive logout requests through the front-channel, look at the [Front-Channel Logout](#_front-channel-logout) client setting. When using this method, consider the following:

- Logout requests sent by Keycloak to clients rely on the browser and on embedded `iframes` that are rendered for the logout page.
- By being based on `iframes`, front-channel logout might be impacted by Content Security Policies (CSP) and logout requests might be blocked.
- If the user closes the browser prior to rendering the logout page or before logout requests are actually sent to clients, their sessions at the client might not be invalidated.

Consider using Back-Channel Logout as it provides a more reliable and secure approach to log out users and terminate their sessions on the clients.

If the client is not enabled with front-channel logout, then Keycloak is going to try first to send logout requests through the back-channel using the [Back-Channel Logout URL](#_back-channel-logout-url). If not defined, the server is going to fall back to using the [Admin URL](#_admin-url).

If even the Admin URL is not specified, Keycloak will not propagate logout to the client. This action can be still sufficient for many client application deployments. For instance, when the client application is a frontend javascript application, it may rely on the [Session Management](#_oidc-logout-session-management) and hence no need exists to send a dedicated logout request to the client application. Other client applications may just rely on the short [Access token lifespan timeout](#_timeouts). This situation means that the session on the application side remains valid for the short time after Keycloak SSO session logout. However, when the client application tries to refresh the access token, the token refresh request to Keycloak will fail because the Keycloak SSO session is already logged out.

##### [](#backchannel-logout)Backchannel Logout

This is a non-browser-based logout that uses direct backchannel communication between Keycloak and clients. Keycloak sends a HTTP POST request containing a logout token to all clients logged into Keycloak. These requests are sent to a registered backchannel logout URLs at Keycloak and are supposed to trigger a logout at client side.

When the client application supports backchannel logout, that logout is usually a better alternative than the front-channel logout as communication happens directly between the client and Keycloak without a web browser involved, which is usually more reliable. Also, when the Keycloak logout of specified user is triggered by the administrator from the [Admin console](#viewing-user-sessions) (or other use of admin REST API) or by the actual user from [the Account console](#_account-service), then Keycloak can propagate a backchannel logout request to the client applications attached to the session being logged out. This scenario is not supported in case of Front-Channel logout because Keycloak can propagate logout to the client by front-channel just when the logout is triggered in the same browser session, which is being logged out. That is not the case for example when session of the specific user is logged-out by the administrator from the admin console.

#### [](#con-server-oidc-uri-endpoints_server_administration_guide)Keycloak server OIDC URI endpoints

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sso-protocols/con-server-oidc-uri-endpoints.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-server-oidc-uri-endpoints.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-server-oidc-uri-endpoints.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The following is a list of OIDC endpoints that Keycloak publishes. These endpoints can be used when a non-Keycloak client adapter uses OIDC to communicate with the authentication server. They are all relative URLs. The root of the URL consists of the HTTP(S) protocol, hostname, and optionally the path: For example

```
https://localhost:8080
```

/realms/{realm-name}/protocol/openid-connect/auth

Used for obtaining a temporary code in the Authorization Code Flow or obtaining tokens using the Implicit Flow, Direct Grants, or Client Grants.

/realms/{realm-name}/protocol/openid-connect/token

Used by the Authorization Code Flow to convert a temporary code into a token.

/realms/{realm-name}/protocol/openid-connect/logout

Used for performing logouts.

/realms/{realm-name}/protocol/openid-connect/userinfo

Used for the User Info service described in the OIDC specification.

/realms/{realm-name}/protocol/openid-connect/revoke

Used for OAuth 2.0 Token Revocation described in [RFC 7009](https://datatracker.ietf.org/doc/html/rfc7009).

/realms/{realm-name}/protocol/openid-connect/certs

Used for the JSON Web Key Set (JWKS) containing the public keys used to verify any JSON Web Token (jwks\_uri)

/realms/{realm-name}/protocol/openid-connect/auth/device

Used for Device Authorization Grant to obtain a device code and a user code.

/realms/{realm-name}/protocol/openid-connect/ext/ciba/auth

This is the URL endpoint for Client Initiated Backchannel Authentication Grant to obtain an auth\_req\_id that identifies the authentication request made by the client.

/realms/{realm-name}/protocol/openid-connect/logout/backchannel-logout

This is the URL endpoint for performing backchannel logouts described in the OIDC specification.

In all of these, replace {realm-name} with the name of the realm.

### [](#_saml)SAML

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sso-protocols/con-saml.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-saml.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-saml.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

[SAML 2.0](https://saml.xml.org/saml-specifications) is a similar specification to OIDC but more mature. It is descended from SOAP and web service messaging specifications so is generally more verbose than OIDC. SAML 2.0 is an authentication protocol that exchanges XML documents between authentication servers and applications. XML signatures and encryption are used to verify requests and responses.

In general, SAML implements two use cases.

The first use case is an application that requests the Keycloak server authenticates a user. Upon successful login, the application will receive an XML document. This document contains an SAML assertion that specifies user attributes. The realm digitally signs the document which contains access information (such as user role mappings) that applications use to determine the resources users are allowed to access in the application.

The second use case is a client accessing remote services. The client requests a SAML assertion from Keycloak to invoke on remote services on behalf of the user.

#### [](#con-saml-bindings_server_administration_guide)SAML bindings

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sso-protocols/con-saml-bindings.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-saml-bindings.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-saml-bindings.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak supports three binding types.

##### [](#redirect-binding)Redirect binding

*Redirect* binding uses a series of browser redirect URIs to exchange information.

1. A user connects to an application using a browser. The application detects the user is not authenticated.
2. The application generates an XML authentication request document and encodes it as a query parameter in a URI. The URI is used to redirect to the Keycloak server. Depending on your settings, the application can also digitally sign the XML document and include the signature as a query parameter in the redirect URI to Keycloak. This signature is used to validate the client that sends the request.
3. The browser redirects to Keycloak.
4. The server extracts the XML auth request document and verifies the digital signature, if required.
5. The user enters their authentication credentials.
6. After authentication, the server generates an XML authentication response document. The document contains a SAML assertion that holds metadata about the user, including name, address, email, and any role mappings the user has. The document is usually digitally signed using XML signatures, and may also be encrypted.
7. The XML authentication response document is encoded as a query parameter in a redirect URI. The URI brings the browser back to the application. The digital signature is also included as a query parameter.
8. The application receives the redirect URI and extracts the XML document.
9. The application verifies the realm’s signature to ensure it is receiving a valid authentication response. The information inside the SAML assertion is used to make access decisions or display user data.

##### [](#post-binding)POST binding

*POST* binding is similar to *Redirect* binding but *POST* binding exchanges XML documents using POST requests instead of using GET requests. *POST* binding uses JavaScript to make the browser send a POST request to the Keycloak server or application when exchanging documents. HTTP responds with an HTML document which contains an HTML form containing embedded JavaScript. When the page loads, the JavaScript automatically invokes the form.

*POST* binding is recommended due to two restrictions:

- **Security** — With *Redirect* binding, the SAML response is part of the URL. It is less secure as it is possible to capture the response in logs.
- **Size** — Sending the document in the HTTP payload provides more scope for large amounts of data than in a limited URL.

##### [](#ecp)ECP

Enhanced Client or Proxy (ECP) is a SAML v.2.0 profile which allows the exchange of SAML attributes outside the context of a web browser. It is often used by REST or SOAP-based clients.

#### [](#keycloak-server-saml-uri-endpoints)Keycloak Server SAML URI Endpoints

Keycloak has one endpoint for all SAML requests.

`http(s)://authserver.host/realms/{realm-name}/protocol/saml`

All bindings use this endpoint.

### [](#ref-saml-vs-oidc_server_administration_guide)OpenID Connect compared to SAML

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sso-protocols/ref-saml-vs-oidc.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fref-saml-vs-oidc.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fref-saml-vs-oidc.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The following lists a number of factors to consider when choosing a protocol.

For most purposes, Keycloak recommends using OIDC.

**OIDC**

- OIDC is specifically designed to work with the web.
- OIDC is suited for HTML5/JavaScript applications because it is easier to implement on the client side than SAML.
- OIDC tokens are in the JSON format which makes them easier for Javascript to consume.
- OIDC has features to make security implementation easier. For example, see the [iframe trick](https://openid.net/specs/openid-connect-session-1_0.html#ChangeNotification) that the specification uses to determine a users login status.

**SAML**

- SAML is designed as a layer to work on top of the web.
- SAML can be more verbose than OIDC.
- Users pick SAML over OIDC because there is a perception that it is mature.
- Users pick SAML over OIDC existing applications that are secured with it.

### [](#_docker)Distribution Registry v2 authentication

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/sso-protocols/con-sso-dist-reg.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-sso-dist-reg.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fsso-protocols%2Fcon-sso-dist-reg.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Docker authentication is disabled by default. To enable docker authentication, see the [Enabling and disabling features](https://www.keycloak.org/server/features) guide.

[Distribution Registry V2 Authentication](https://distribution.github.io/distribution/spec/auth/) is a protocol, similar to OIDC, that authenticates users against Distribution registries. Keycloak’s implementation of this protocol lets Docker clients use a Keycloak authentication server authenticate against a registry. This protocol uses standard token and signature mechanisms but it does deviate from a true OIDC implementation. It deviates by using a very specific JSON format for requests and responses as well as mapping repository names and permissions to the OAuth scope mechanism.

#### [](#docker-authentication-flow)Docker authentication flow

The authentication flow is described in the [Docker API documentation](https://distribution.github.io/distribution/spec/auth/token/). The following is a summary from the perspective of the Keycloak authentication server:

- Perform a `docker login`.
- The Docker client requests a resource from the Distribution registry. If the resource is protected and no authentication token is in the request, the Distribution registry server responds with a 401 HTTP message with some information on the permissions that are required and the location of the authorization server.
- The Docker client constructs an authentication request based on the 401 HTTP message from the Distribution registry. The client uses the locally cached credentials (from the `docker login` command) as part of the [HTTP Basic Authentication](https://datatracker.ietf.org/doc/html/rfc2617) request to the Keycloak authentication server.
- The Keycloak authentication server attempts to authenticate the user and return a JSON body containing an OAuth-style Bearer token.
- The Docker client receives a bearer token from the JSON response and uses it in the authorization header to request the protected resource.
- The Distribution registry receives the new request for the protected resource with the token from the Keycloak server. The registry validates the token and grants access to the requested resource (if appropriate).

Keycloak does not create a browser SSO session after successful authentication with the Docker protocol. The browser SSO session does not use the Docker protocol as it cannot refresh tokens or obtain the status of a token or session from the Keycloak server; therefore a browser SSO session is not necessary. For more details, see the [transient session](#_transient-session) section.

#### [](#keycloak-distribution-registry-v2-authentication-server-uri-endpoints)Keycloak Distribution registry v2 Authentication Server URI Endpoints

Keycloak has one endpoint for all Docker auth v2 requests.

`http(s)://authserver.host/realms/{realm-name}/protocol/docker-v2/auth`

## [](#_admin_permissions)Managing access to realm resources

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/admin-console-permissions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fadmin-console-permissions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fadmin-console-permissions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Each realm created on the Keycloak has a dedicated Admin Console from which that realm can be managed. The `master` realm is a special realm that allows admins to manage more than one realm on the system. You can also define fine-grained access to users in different realms to manage the server. This chapter goes over all the scenarios for this.

### [](#_master_realm_access_control)Master realm access control

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/admin-console-permissions/master-realm.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fadmin-console-permissions%2Fmaster-realm.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fadmin-console-permissions%2Fmaster-realm.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The `master` realm in Keycloak is a special realm and treated differently than other realms. Users in the Keycloak `master` realm can be granted permission to manage zero or more realms that are deployed on the Keycloak server. When a realm is created, Keycloak automatically creates various roles that grant permissions to access that new realm. Access to The Admin Console and Admin REST endpoints can be controlled by mapping these roles to users in the `master` realm. It’s possible to create multiple superusers, as well as users that can only manage specific realms.

#### [](#global-roles)Global roles

There are two realm-level roles in the `master` realm. These are:

- admin
- create-realm

Users with the `admin` role are superusers and have full access to manage any realm on the server. Users with the `create-realm` role are allowed to create new realms. They will be granted full access to any new realm they create.

Admins with the `create-realm` role who do not have the `admin` role will also need the realm-specific `query-realms` role to be able to create and list realms in the Admin Console.

#### [](#realm-specific-roles)Realm specific roles

Admin users within the `master` realm can be granted management privileges to one or more other realms in the system. Each realm in Keycloak is represented by a client in the `master` realm. The name of the client is `<realm name>-realm`. These clients each have client-level roles defined which define varying level of access to manage an individual realm.

The roles available are:

- create-client
- impersonation
- manage-authorization
- manage-clients
- manage-events
- manage-identity-providers
- manage-realm
- manage-users
- query-clients
- query-groups
- query-realms
- query-users
- view-authorization
- view-clients
- view-events
- view-identity-providers
- view-realm
- view-users

Assign the roles you want to your users and they will only be able to use that specific part of the administration console.

Admins with the `manage-users` role will only be able to assign admin roles to users that they themselves have. So, if an admin has the `manage-users` role but doesn’t have the `manage-realm` role, they will not be able to assign this role.

### [](#_per_realm_admin_permissions)Dedicated realm admin consoles

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/admin-console-permissions/per-realm.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fadmin-console-permissions%2Fper-realm.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fadmin-console-permissions%2Fper-realm.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Each realm has a dedicated Admin Console that can be accessed by going to the url `/admin/{realm-name}/console`. Users within that realm can be granted realm management permissions by assigning specific user role mappings.

Each realm has a built-in client called `realm-management`. You can view this client by going to the `Clients` left menu item of your realm. This client defines client-level roles that specify permissions that can be granted to manage the realm.

- create-client
- impersonation
- manage-authorization
- manage-clients
- manage-events
- manage-identity-providers
- manage-realm
- manage-users
- query-clients
- query-groups
- query-realms
- query-users
- realm-admin
- view-authorization
- view-clients
- view-events
- view-identity-providers
- view-realm
- view-users

Assign the roles you want to your users and they will only be able to use that specific part of the administration console.

### [](#_fine_grained_permissions)Delegating realm administration using permissions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/admin-console-permissions/fine-grain-v2.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fadmin-console-permissions%2Ffine-grain-v2.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fadmin-console-permissions%2Ffine-grain-v2.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can delegate realm management to other administrators, the realm administrators, using the fine-grained admin permissions feature. Different from the Role-Based Access Control (RBAC) Mechanism provided through the [Global and Realm specific roles](#_master_realm_access_control), this feature provides a more fine-grained control over how realm resources can be accessed and managed based on a well-defined set of operations that can be performed on them.

By relying on a Policy-Based Access Control, realm administrators can define permissions to realm resources such as users, groups, and clients, using different policy types, or access control methods, so that a realm administrator is limited to access a subset of realm resources and their operations.

The feature provides an alternative to the aforementioned RBAC mechanism, but it does not replace it. You are still able to grant administrative roles like `view-users` or `manage-clients` to delegate access to realm administrators but doing so will skip the mechanisms provided by this feature.

Enforcing access to realm resources only applies when managing resources through the administration console or the Admin API.

#### [](#_understanding_different_types_realm_admins_)Understanding the different types of realm administrators

There are three different types of realm administrators when managing a realm:

- **Server administrators**: These are users (or service accounts) that have been granted with the `admin` role in the master realm.
- **Realm administrators**: These are users (or service accounts) that have been granted with the `realm-admin` role in a specific realm.
- **Delegated realm administrators**: These are users (or service accounts) other than the types above but granted access to manage a realm through the fine-grained admin permissions feature.

Both server and realm administrators can manage a realm with full access to all their resources and administrative operations. The main difference between them is that server administrators can manage any realm in the server, while realm administrators can only manage the realm they have been granted access to.

Delegated realm administrators, on the other hand, can have limited access to a realm based on the permissions defined through this feature. They can only access the resources and operations they have been granted access to.

Be aware that both server and realm administrators are not affected by the permissions you define when managing access to realm resources through this feature. Always make sure to review the users granted with the `admin` or `realm-admin` roles to avoid any potential privilege escalation.

#### [](#understanding-the-realm-resource-types)Understanding the Realm Resource Types

In a realm, you can manage different types of resources such as users, groups, clients, client scopes, roles, and so on. As a realm administrator, you are constantly managing these resources when managing identities and how they authenticate and are authorized to access a realm and applications.

This feature provides the necessary mechanisms to enforce access controls when managing realm resources, limited to:

- Users
- Groups
- Clients
- Roles

You can manage permissions for all resources of a given resource type, such as all users in a realm, or for a specific realm resource, such as a specific user or set of users in the realm.

#### [](#understanding-the-scopes-of-access)Understanding the scopes of access

Each realm resource supports a well-defined set of management operations, or scopes, that can be performed on them, such as `view`, `manage`, and resource-specific operations such as `view-members`, if you take groups as an example.

When managing permissions, you are selecting a set of one or more scopes from a resource type to allow realm administrators to perform specific operations on a resource type. For instance, granting a `view` scope will give access to realm administrators to list, search, and view a realm resource. On the other hand, the `manage` scope will allow administrators to perform updates and deletes on them.

The scopes are completely independent of each other. If you give access to `manage` a realm resource, that does not mean the `view` scope is granted automatically. No transitive dependency exists between scopes. Although this might impact the overall user experience when managing permissions because you need to select individual scopes, the benefit is that you can more easily identify the permissions that enforce access to a specific scope.

Certain scopes from a resource type have a relationship (not a transitive dependency) to scopes in another resource type. This relationship is mainly true when you manage a resource type that represents a group of realm resources, such as realm groups and their members.

##### [](#users-resource-type)Users Resource Type

The **Users** realm resource type represents the users in a realm. You can manage permissions for users based on the following set of scopes:

   **Scope** **Description** **Also granted by**

**view**

Defines if a realm administrator can view users. This scope should be set whenever you want to make users available from queries.

`view-members`

**manage**

Defines if a realm administrator can manage users.

`manage-members`

**manage-group-membership**

Defines if a realm administrator can assign or unassign users to/from groups.

`manage-membership-of-members`

**map-roles**

Defines if a realm administrator can assign or unassign roles to/from users.

None

**impersonate**

Defines if a realm administrator can impersonate other users.

`impersonate-members`

**reset-password**

Defines if a realm administrator can reset user passwords. If no permission with `reset-password` is found, it falls back to check `manage` scope.

None

The user resource type has a strong relationship with some of the permissions you can set to groups. Most of the time, users are members of groups and granting access to `view-members` or `manage-members` of a group should also allow a realm administrator to `view` and `manage` members of that group.

To assign or unassign a specific user to or from a group, both sides of the relationship are required: `manage-group-membership` on the user and `manage-membership` on the target group. The `manage-membership-of-members` scope can be used to grant or deny `manage-group-membership` based on the groups that the target user currently belongs to.

This feature does not support enforcing access to federated resource, however, this limitation is being considered for future improvement.

##### [](#groups-resource-type)Groups Resource Type

The **Groups** realm resource type represents the groups in a realm. You can manage permissions for groups based on the following set of management operations:

  **Operation** **Description**

**view**

Defines if a realm administrator can view groups. This scope should be set whenever you want to make groups available from queries.

**manage**

Defines if a realm administrator can manage groups.

**view-members**

Defines if a realm administrator can view group members. This operation applies to any child group in the group hierarchy. This can be prevented by explicitly denying permission for specific subgroups.

**manage-members**

Defines if a realm administrator can manage group members. This operation applies to any child group in the group hierarchy. This can be prevented by explicitly denying permission for specific subgroups.

**impersonate-members**

Defines if a realm administrator can impersonate group members. This operation applies to any child group in the group hierarchy. This can be prevented by explicitly denying permission for specific subgroups.

**manage-membership**

Defines if a realm administrator can add or remove members from groups.

**manage-membership-of-members**

Defines if a realm administrator can grant or deny `manage-group-membership` for members of a group.

##### [](#clients-resource-type)Clients Resource Type

The **Clients** realm resource type represents the clients in a realm. You can manage permissions for clients based on the following set of management operations:

  **Operation** **Description**

**view**

Defines if a realm administrator can view clients. This scope should be set whenever you want to make clients available from queries.

**manage**

Defines if a realm administrator can manage clients.

**map-roles**

Defines if a realm administrator can assign any role defined by a client to a user.

**map-roles-composite**

Defines if a realm administrator can assign any role defined by a client as a composite to another role.

**map-roles-client-scope**

Define if a realm administrator can assign any role defined by a client to a client scope.

The **map-roles** operation does not grant the ability to manage users or assign roles arbitrarily. The administrator must also have user role mapping permissions on the user.

##### [](#roles-resource-type)Roles Resource Type

The **Roles** realm resource type represents the roles in a realm. You can manage permissions for roles based on the following set of management operations:

  **Operation** **Description**

**map-role**

Defines if a realm administrator can assign a role (or multiple roles) to a user.

**map-role-composite**

Defines if a realm administrator can assign a role (or multiple roles) as a composite to another role.

**map-role-client-scope**

Defines if a realm administrator can apply a role (or multiple roles) to a client scope.

The **map-roles** operation does not grant the ability to manage users or assign roles arbitrarily. The administrator must also have user role mapping permissions on the user.

If there is a client resource type permission for the **map-roles**, **map-roles-composite**, or **map-roles-client-scope** scopes, it will take precedence over any role resource type permission if the role is a client role.

#### [](#enabling-admin-permissions-to-a-realm)Enabling admin permissions to a realm

To enable fine-grained admin permissions in a realm, follow these steps:

- Log in to the Admin Console.
- Click **Realm settings**.
- Enable **Admin Permissions** and click **Save**.

![Fine grain enable](./images/fine-grain-enable.png)

Once enabled, a **Permissions** section appears in the left-side menu of the administration console.

![Fine grain permissions tab](./images/fine-grain-permissions-tab.png)

From this section, you can manage the permissions for realm resources.

#### [](#_managing-permissions)Managing Permissions

The **Permissions** tab provides an overview of all active permissions within a realm. From here, administrators can create, update, delete, or search for permissions. You can also pre-evaluate the permissions you have created to check if they are enforcing access to realm resources as expected. For more details, see [Evaluating Permissions](#_evaluating-permissions).

To create a permission, click on the `Create permission` button and select the resource type you want to protect.

![Selecting a resource type to protect](./images/select-resource-type.png)

Once you select the resource type, you can now define how access should be enforced for a set of one or more resources of the selected type:

![Creating a permission](./images/create-permission.png)

When managing a permission you can define the following settings:

- **Name**: A unique name for the permission. The name should also not conflict with any policy name
- **Description**: An optional description to better describe what the permission is about
- **Authorization scopes**: A set of one or more scopes representing the operations you want to protect for the selected resource type. An administrator must have explicit permission assigned for each operation to perform the corresponding action. For example, assigning only **manage** without **view** will prevent the user from being visible.
- **Enforce access to**: Defines if the permission should enforce access to all resources of the selected type or to specific resources in a realm.
- **Policies**: Defines a set of one or more policies that should be evaluated to grant or deny access to the selected resource(s).

After creating the permission, it will automatically take effect when enforcing access to (all) resources and scopes you selected. Keep that fact in mind when creating and updating permissions in production.

##### [](#defining-permissions-for-viewing-realm-resources)Defining permissions for viewing realm resources

This feature relies on a partial evaluation mechanism to partially evaluate the permissions that a realm administrator has when listing and viewing realm resources. This mechanism will pre-fetch all the permissions set for view-related scopes where the realm administrator is referenced either directly or indirectly.

Permissions that grant access to `view` a realm resource of a certain type must use one of the following policies to make them available from queries:

- `User`
- `Group`
- `Role`
- `Aggregated`

In case of using an `Aggregated` policy, all the underlying policies must be of type `User`, `Group`, or `Role`. Otherwise, the policy will always evaluate to `DENY` when performing partial evaluation.

By using any of the policies above, Keycloak can pre-calculate the set of resources that a realm administration can view by looking for a direct (if using a user policy) or indirect (if using a role or group policy) reference to the realm administrator. Therefore, the partial evaluation mechanism involves decorating queries with access controls that will run at the database level. This capability is mainly important to properly allow paginating resources as well as avoid an additional overhead on the server-side when evaluating permissions for each realm resource returned by queries.

Partial evaluation and filtering occurs only if the feature is enabled to a realm, and if the user is not granted with view-related administrative roles like `view-users` or `view-clients`. For instance, it will not happen for administrators granted with the `admin` role at the master realm (server administrators), or realm administrators granted with the `realm-admin` role in a realm other than the master realm.

When querying resources, the partial evaluation mechanism works as follows:

- Resolve all the permissions for a certain resource type that reference the realm administrator
- Pre-evaluate each permission to check if the realm administrator does or does not have access to the resources associated with the permission
- Decorate database queries based on the resources granted or denied

As a result, the result set of a query will hold only the realm resources where realm administrators have access to any of the view-related scopes.

##### [](#searching-permissions)Searching Permissions

The Admin Console provides several ways to search for permissions, supporting the following capabilities:

- Search for permissions that contain a specific string in their **Name**
- Search for permissions of a specific resource type, such as **Users**
- Search for permissions of a specific resource type that apply to a particular resource (such as **Users** resource type for user `myadmin`).
- Search for permissions of a specific resource type with a given scope (such as **Users** resource type permissions with the **manage** scope).
- Search for permissions of a specific resource type that apply to a particular resource and have a specific scope (such as **Users** resource type permissions with the **manage** scope for user `myadmin`).

Fine grained permissions search

![Fine grained permissions search](./images/fine-grain-search.png)

These capabilities allow realm administrators to perform queries on their universe of permissions and identify which ones are enforcing access to a set of one or more realm resources and their scopes. Combined with the evaluation tool on the **Evaluation** tab, they provide a key management tool for managing permissions in a realm. See [Evaluating Permissions](#_evaluating-permissions) for more details.

#### [](#managing-policies)Managing Policies

The **Policies** tab allows administrators to define conditions using different access control methods to determine whether a permission should be granted to an administrator attempting to perform operations on a realm resource. When managing permissions, you must associate at least a single policy to grant or deny access to a realm resource.

Policies are basically conditions that will evaluate to either a `GRANT` or a `DENY`. Their outcome will decide whether a permission should be granted or denied.

A permission is only granted if all its associated policies evaluate to a `GRANT`. Otherwise, the permission is denied and a realm administrator will not be able to access the protected resource.

Keycloak provides a set of built-in policies that you can choose from:

![Selecting a policy type](./images/select-policy-type.png)

Once you have a well-defined and stable permission model for your realm, less need exists to create policies. You can instead reuse existing policies to create more permissions.

For more details about each policy type, see [Managing policies](https://www.keycloak.org/docs/26.6.3/authorization_services/#_policy_overview).

#### [](#_evaluating-permissions)Evaluating Permissions

The **Evaluation** tab provides a testing environment where administrators can verify that permissions are enforcing access as expected. The administrator can see what permissions are involved when enforcing access to a particular resource and what the outcome is.

You need to provide a set of fields in order to run an evaluation:

- `User`, the realm administrator or the subject trying to access a resource
- `Resource Type`, the resource type you want to evaluate
- `Resource Selector`, depending on the selected `Resource Type` you will be prompted to select a specific realm resource like a user, group, or client.
- `Authorization scope`, the scope or the operation you want to evaluate. If not provided, the evaluation will happen for all the scopes of the selected resource type.

Fine grained permissions evaluation tab

![Fine grained permissions evaluation tab](./images/fine-grain-evaluation.png)

By clicking the `Evaluate` button, the server will evaluate all the permissions associated with the selected resource and scopes just like if the selected `User` were trying to access the resource when using the administration console or the Admin API.

For instance, in the example above you can see that the user `myadmin` can **manage** user `user-1` because a `Allow managing all realm users` permission voted to a `PERMIT`, therefore granting access to the `manage` scope. However, all the other scopes were denied.

Combined with the searching capabilities from the **Permissions** tab, you can perform troubleshooting to identify any permission that is not behaving as expected.

When evaluating permissions, the following rules apply:

- The outcome from resource-specific permissions have precedence over broader permissions that give access to all resources of a certain type
- If no permissions exist for a specific resource, access will be granted based on the permission that grants access to all resources of a certain type
- The outcome from different permissions that enforce access to a specific resource will only grant access if they all permit access to the resource

##### [](#_resolving-conflicting-permissions)Resolving conflicting permissions

Permissions can have multiple policies associated with them. As the authorization model evolves, it is common for some policies within a permission or even different permissions related to a specific resource to conflict.

The evaluation outcome will be "denied" whenever any permission is evaluated to "DENY." If there are multiple permissions related to the same resource, all of them must grant access in order for the outcome to be "granted."

Fine-grained admin permissions allow you to set up permissions for individual resources or for the resource type itself (such as all users, all groups, and so on.). If a permission or permissions related to a specific resource exist, the "all-resource" permission is **NOT** taken into account during evaluation. If no specific permission exists, the fallback is to the "all-resource" permission. This approach helps address scenarios like allowing members of the `realm-admins` group to manage members of realm groups, but preventing them from managing members of the `realm-admins` group themselves.

#### [](#_realm_access_control)Accessing a Realm administration console as a Realm Administrator

Realm administrators can access a dedicated realm-specific administration console that allows them to manage resources within their assigned realm. This console is separate from the main Keycloak Admin Console, which is typically used by realm administrators.

For more details on dedicated realm administration consoles and available roles, refer to: [Dedicated admin consoles](#_per_realm_admin_permissions).

To access the administration console, a realm administrator must have at least one of the following roles assigned, depending on the resources they need to administer:

- **query-users** – Required to query realm users.
- **query-groups** – Required to query realm groups.
- **query-clients** – Required to query realm clients.

By granting any of these roles to a realm user, they will be able to access the administration console, but only for the areas that correspond to roles granted. For instance, if you assign the `query-users` role, the realm administrator will only have access to the `Users` section in the administration console. If an administrator is responsible for multiple resource types (such as both users and groups), they must have all the corresponding "query-\*" roles assigned.

These roles enable basic access to query resources but do not grant permission to view or modify them. To grant or deny access to realm resources you need to set up the permissions for any of the operations available from each resource type. For more details, see [Managing Permissions](#_managing-permissions).

##### [](#roles-and-permission-relationship)Roles and Permission relationship

Fine grained permissions are used to grant additional permissions. You cannot override the default behavior of the built-in admin roles. If a realm administrator is assigned one or more admin roles, it prevents the permissions from being evaluated. This means that if a respective admin role is assigned to a realm administrator, permission evaluation will be bypassed, and access will be granted.

  **Admin Role** **Description**

**query-users**

A realm administrator can see the **Users** section in administration console and can search for users in the realm. It does not grant the ability to **view** users.

**query-groups**

A realm administrator can see the **Groups** section in administration console and can search for groups in the realm. It does not grant the ability to **view** groups.

**query-clients**

A realm administrator can see the **Clients** section in administration console and can search for clients in the realm. It does not grant the ability to **view** clients.

**view-users**

A realm administrator can **view** all users and groups in the realm.

**manage-users**

A realm administrator can **view**, **map-roles**, **manage-group-membership** and **manage** all users in the realm, as well as **view**, **manage-membership**, **manage-membership-of-members** and **manage** groups in the realm.

**impersonation**

A realm administrator can **impersonate** all users in the realm.

**view-clients**

A realm administrator can **view** all clients in the realm.

**manage-clients**

A realm administrator can **view** and **manage** all clients and client scopes in the realm.

When this feature is enabled in a realm, only server and realm administrators with the corresponding admin roles can grant these roles to other realm administrators. Delegated realm administrators cannot assign administrative roles to other realm administrators.

#### [](#understanding-some-common-use-cases)Understanding some common use cases

Consider a situation where an administrator wants to allow a group of administrators to manage all users in the realm except those that belong to the administrators group. This example includes a `test` realm and a `test-admins` group.

##### [](#allowing-to-manage-users-by-group-of-administrators)Allowing to manage users by group of administrators

Create user permission permission, allowing to view and manage all users in the realm for members of the `test-admins` group:

- Navigate to the **Permissions** tab in the administration console.
- Click **Create permission** and choose **Users** resource type.
- Fill in the name, such as `Disallow managing test-admins`.
- Choose **view** and **manage** authorization scopes, keep checked **All Users**.
- Create a condition, which needs to be met to get an access by clicking **Create new policy**.
- Fill in the name `Allow test-admins`, select **Group** as **Policy type**.
- Click **Add groups** button and select `test-admins` group, click **Save**.
- Click **Save** on **Create permission** page.

##### [](#allowing-to-manage-users-by-group-of-admins-but-not-group-members)Allowing to manage users by group of admins but not group members

Let’s exlude the members of the group itself, so that `test-admins` cannot manage other admins.

- Create new permission by clicking **Create permission**.
- This time choose **Groups** resource type.
- Fill in the name, such as `Disallow managing test-admins`.
- Choose **manage-members** authorization scope.
- Select **Specific Groups** and choose `test-admins` group.
- **Create new policy** of type **Group**.
- Fill the name `Disallow test-admins` and select `test-admins` group.
- Switch to **Negative Logic** for the policy, **Save** the policy
- **Save** the permission

##### [](#allowing-to-impersonate-users-for-members-of-a-group-with-a-specific-role-assigned)Allowing to impersonate users for members of a group with a specific role assigned

- Create a "User Permission" for specific users (or all users) you want to allow impersonation.
- Create a "Group Policy" allowing access to members of `test-admins`.
- Create a "Role Policy" allowing access to users assigned the `impersonation-admin` role.
- Assign both policies to the permission.

##### [](#preventing-specific-users-from-being-impersonated)Preventing specific users from being impersonated

- Create a **User Permission** for the specific users you want to prevent from being impersonated.
- Create any policy that evaluates to deny (such as a user policy with no users selected).
- Assign the policy to the permission to effectively block impersonation for the selected users.

##### [](#allowing-to-view-users-but-not-managing-them-for-admins-with-a-defined-role-assigned)Allowing to view users but not managing them for admins with a defined role assigned

- Create a "User Permission" with the **view** scope for all users.
- Create a "Role Policy" allowing access to users with specific role assigned.
- Do *not* assign the `manage` scope to prevent modification of user details.

##### [](#allowing-to-manage-users-and-role-assignment-for-members-of-a-group)Allowing to manage users and role assignment for members of a group

- Create a "User Permission" with the **manage**, **map-roles** scopes for all users.
- Create a "Group Policy" allowing access to members of `test-admins`.

##### [](#allowing-to-view-and-manage-members-of-a-group-but-not-members-of-its-subgroups)Allowing to view and manage members of a group but not members of its subgroups

- Create a "Group Permission" with the **view-members** and **manage-members** scopes for specific group `mygroup`.
- Assign a "Group Policy" targeting `test-admins` to it.
- Create another "Group Permission" with the **view-members** and **manage-members** scopes for specific group, select all subgroups of the `mygroup`.
- Create negative "Group Policy" for `test-admins` and assign it to the "subgroups" permission.

##### [](#allowing-to-impersonate-members-of-a-specific-group)Allowing to impersonate members of a specific group

- Create a "Group Permission" with the **impersonate-members** for specific group `mygroup`.
- Assign a "Group Policy" targeting `mygroup-helpdesk` to it.

#### [](#performance-considerations)Performance considerations

When enabling the feature to a realm, there is an additional overhead when realm administrators are managing any of the supported resource types. This is mainly true when performing these operations:

- Listing and searching
- Updating or deleting

The feature introduces additional checks whenever you are listing or managing realm resources in order to enforce access based on the permissions you have defined. This is mainly true when querying realm resources due to the additional overhead to partially evaluate the permissions for a realm administrator to filter and paginate the results.

Fewer permissions referencing a realm administrator user and most of the resources they can access is better. For instance, if you want to delegate access to a realm administrator to manage users, it is better to have those users as members of a group. By doing that, you are improving not only the performance when evaluating permissions but also creating a permission model that is easier to manage.

The main impact of access enforcement is when querying realm resources. If a realm administrator is, for instance, referenced in thousands of permissions through a user, role, or group policy, the partial evaluation mechanism that happens when querying realm resources will query all those permissions from the database. A more concise and optimized model will help to fetch fewer permissions but the enough to grant or deny access to realm resources.

For instance, granting access to a realm administrator to view and manage users in a realm is better done with a group permission than create individual permissions for each individual user in a realm. As well as make sure the policies associated with a permission referencing a realm administrator either by a direct reference (user policy), or indirect (role or group policy) reference, do not span multiple (thousands of) permissions, regardless of the resource type.

As an example, suppose you have three users in a realm, and you want to allow `bob`, a realm administrator, to `view` and `manage` them. A non-optimal permission model would create three different permissions, for each user, where a user policy grants access to `bob`. Instead, you can have a single group permission, or even a single user permission, that groups those three users while still granting access to `bob` using the same user policy.

The same is true if you want to give access to more realm administrators to those three users. Instead of creating individual policies, you can consider using a group or role policy instead. The permission model is use-case-specific, but these recommendations are important to provide not only better manageability but also improve the overall performance of the server when managing realm resources.

In terms of server configuration, depending on the size of your realm and the number of permissions and policies you have, you might consider changing the cache configuration to increase the size of the following caches:

- `realms`
- `users`
- `authorization`

Consider looking at the server metrics for these caches to find the best value when sizing your deployment.

When filtering resources, the partial evaluation mechanism will eventually rely on `IN` clauses in SQL statements to filter the results. Depending on your database, you might have limitations on the number of parameters for the `IN` clause. That is the case for old versions of the Oracle database, which has a hard limit to 1000 parameters. To avoid such problems, keep in mind the considerations above about the number of permissions that grants or deny access to a single realm administrator.

### [](#fine-grained-admin-permissions-v1)Fine grained admin permissions V1

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/admin-console-permissions/fine-grain.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fadmin-console-permissions%2Ffine-grain.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fadmin-console-permissions%2Ffine-grain.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Preview feature fine-grained admin permissions V1 has been replaced by a new [supported version](#_fine_grained_permissions). Version 1 of the feature is still available, but it has been deprecated and will be removed in future release. To enable it, start the server with `--features=admin-fine-grained-authz:v1`.

Sometimes roles like `manage-realm` or `manage-users` are too coarse grain and you want to create restricted admin accounts that have more fine grain permissions. Keycloak allows you to define and assign restricted access policies for managing a realm. Things like:

- Managing one specific client
- Managing users that belong to a specific group
- Managing membership of a group
- Limited user management.
- Fine grain impersonation control
- Being able to assign a specific restricted set of roles to users.
- Being able to assign a specific restricted set of roles to a composite role.
- Being able to assign a specific restricted set of roles to a client’s scope.
- New general policies for viewing and managing users, groups, roles, and clients.

There are some important things to note about fine grain admin permissions:

- Fine grain admin permissions were implemented on top of [Authorization Services](https://www.keycloak.org/docs/26.6.3/authorization_services/). It is highly recommended that you read up on those features before diving into fine grain permissions.
- Fine grain permissions are only available within [dedicated admin consoles](#_per_realm_admin_permissions) and admins defined within those realms. You cannot define cross-realm fine grain permissions.
- Fine grain permissions are used to grant additional permissions. You cannot override the default behavior of the built-in admin roles.

#### [](#managing-one-specific-client)Managing one specific client

Let’s look first at allowing an admin to manage one client and one client only. In our example, we have a realm called `test` and a client called `sales-application`. In the realm `test` we will give a user in that realm permission to only manage that application.

You cannot do cross realm fine grain permissions. Admins in the `master` realm are limited to the predefined admin roles defined in previous chapters.

##### [](#permission-setup)Permission setup

The first thing we must do is login to the Admin Console so we can set up permissions for that client. We navigate to the management section of the client, we want to define fine-grain permissions for.

Client management

![Fine grain client](./images/fine-grain-client.png)

You should see a tab menu item called `Permissions`. Click on that tab.

Client permissions tab

![Fine grain client permissions tab](./images/fine-grain-client-permissions-tab-off.png)

By default, each client is not enabled to do fine grain permissions. So turn the `Permissions Enabled` switch to on to initialize permissions.

If you turn the `Permissions Enabled` switch to off, it will delete any and all permissions you have defined for this client.

Client permissions tab

![Fine grain permission tab](./images/fine-grain-client-permissions-tab-on.png)

When you switch `Permissions Enabled` to on, it initializes various permission objects behind the scenes using [Authorization Services](https://www.keycloak.org/docs/26.6.3/authorization_services/). For this example, we’re interested in the `manage` permission for the client. Clicking on that will redirect you to the permission that handles the `manage` permission for the client. All authorization objects are contained in the `realm-management` client’s `Authorization` tab.

Client manage permission

![Fine grain client manage permission](./images/fine-grain-client-manage-permissions.png)

When first initialized the `manage` permission does not have any policies associated with it. You will need to create one by going to the policy tab. To get there fast, click on the `Client details` link shown in the above image. Then click on the policies tab.

On this page, look for the `Create client policy` button, which you can use to define many policies. You can define a policy that is associated with a role or a group or even define rules in JavaScript. For this simple example, we are going to create a `User Policy`.

User policy

![Fine grain client user policy](./images/fine-grain-client-user-policy.png)

This policy will match a hard-coded user in the user database. In this case, it is the `sales-admin` user. We must then go back to the `sales-application` client’s `manage` permission page and assign the policy to the permission object.

Assign user policy

![Fine grain client assign user policy](./images/fine-grain-client-assign-user-policy.png)

The `sales-admin` user now has permission to manage the `sales-application` client.

There is one more thing we have to do. Go to `Users`, select the `sales-admin` user, then go to the `Role Mappings` tab and assign the `query-clients` role to the user.

Assign query-clients

![Fine grain assign query clients](./images/fine-grain-assign-query-clients.png)

Why do you have to do this? This role tells the Admin Console what menu items to render when the `sales-admin` visits the Admin Console. The `query-clients` role tells the Admin Console that it should render client menus for the `sales-admin` user.

IMPORTANT If you do not set the `query-clients` role, restricted admins like `sales-admin` will not see any menu options when they log into the Admin Console

##### [](#testing-it-out)Testing it out

Next, we log out of the master realm and re-login to the [dedicated admin console](#_per_realm_admin_permissions) for the `test` realm using the `sales-admin` as a username. This is located under `/admin/test/console`.

Sales admin login

![Fine grain sales admin login](./images/fine-grain-sales-admin-login.png)

This admin is now able to manage this one client.

#### [](#restrict-user-role-mapping)Restrict user role mapping

Another thing you might want to do is to restrict the set of roles an admin is allowed to assign to a user. Continuing our last example, let’s expand the permission set of the 'sales-admin' user so that he can also control which users are allowed to access this application. Through fine grain permissions, we can enable it so that the `sales-admin` can only assign roles that grant specific access to the `sales-application`. We can also restrict it so that the admin can only map roles and not perform any other types of user administration.

The `sales-application` has defined three different client roles.

Sales application roles

![Fine grain sales application roles](./images/fine-grain-sales-application-roles.png)

We want the `sales-admin` user to be able to map these roles to any user in the system. The first step to do this is to allow the role to be mapped by the admin. If we click on the `viewLeads` role, you’ll see that there is a `Permissions` tab for this role.

View leads role permission tab

![Fine grain view leads role](./images/fine-grain-view-leads-role-tab.png)

If we click on that tab and turn the `Permissions Enabled` on, you’ll see that there are a number of actions we can apply policies to.

View leads permissions

![Fine grain view leads permissions](./images/fine-grain-view-leads-permissions.png)

The one we are interested in is `map-role`. Click on this permission and add the same User Policy that was created in the earlier example.

Map-roles permission

![Fine grain map roles permission](./images/fine-grain-map-roles-permission.png)

What we’ve done is say that the `sales-admin` can map the `viewLeads` role. What we have not done is specify which users the admin is allowed to map this role too. To do that we must go to the `Users` section of the admin console for this realm. Clicking on the `Users` left menu item brings us to the users interface of the realm. You should see a `Permissions` tab. Click on that and enable it.

Users permissions

![Fine grain user permissions](./images/fine-grain-users-permissions.png)

The permission we are interested in is `map-roles`. This is a restrictive policy in that it only allows admins the ability to map roles to a user. If we click on the `map-roles` permission and again add the User Policy we created for this, our `sales-admin` will be able to map roles to any user.

The last thing we have to do is add the `view-users` role to the `sales-admin`. This will allow the admin to view users in the realm he wants to add the `sales-application` roles to.

Add view-users

![Fine grain add view users](./images/fine-grain-add-view-users.png)

##### [](#testing-it-out-2)Testing it out

Next, we log out of the master realm and re-login to the [dedicated admin console](#_per_realm_admin_permissions) for the `test` realm using the `sales-admin` as a username. This is located under `/admin/test/console`.

You will see that now the `sales-admin` can view users in the system. If you select one of the users you’ll see that each user detail page is read only, except for the `Role Mappings` tab. Going to this tab you’ll find that there are no `Available` roles for the admin to map to the user except when we browse the `sales-application` roles.

Assign viewLeads

![Fine grain add view leads](./images/fine-grain-add-view-leads.png)

We’ve only specified that the `sales-admin` can map the `viewLeads` role.

##### [](#per-client-map-roles-shortcut)Per client map-roles shortcut

It would be tedious if we had to do this for every client role that the `sales-application` published. to make things easier, there’s a way to specify that an admin can map any role defined by a client. If we log back into the admin console to our master realm admin and go back to the `sales-application` permissions page, you’ll see the `map-roles` permission.

Client map-roles permission

![Fine grain client permissions](./images/fine-grain-client-permissions-tab-on.png)

If you grant access to this particular permission to an admin, that admin will be able map any role defined by the client.

#### [](#full-list-of-permissions)Full list of permissions

You can do a lot more with fine grain permissions beyond managing a specific client or the specific roles of a client. This chapter defines the whole list of permission types that can be described for a realm.

##### [](#role)Role

When going to the `Permissions` tab for a specific role, you will see these permission types listed.

map-role

Policies that decide if an admin can map this role to a user. These policies only specify that the role can be mapped to a user, not that the admin is allowed to perform user role mapping tasks. The admin will also have to have manage or role mapping permissions. See [Users Permissions](#_users-permissions) for more information.

map-role-composite

Policies that decide if an admin can map this role as a composite to another role. An admin can define roles for a client if he has to manage permissions for that client but he will not be able to add composites to those roles unless he has the `map-role-composite` privileges for the role he wants to add as a composite.

map-role-client-scope

Policies that decide if an admin can apply this role to the scope of a client. Even if the admin can manage the client, he will not have permission to create tokens for that client that contain this role unless this privilege is granted.

##### [](#client)Client

When going to the `Permissions` tab for a specific client, you will see these permission types listed.

view

Policies that decide if an admin can view the client’s configuration.

manage

Policies that decide if an admin can view and manage the client’s configuration. There are some issues with this in that privileges could be leaked unintentionally. For example, the admin could define a protocol mapper that hardcoded a role even if the admin does not have privileges to map the role to the client’s scope. This is currently the limitation of protocol mappers as they don’t have a way to assign individual permissions to them like roles do.

configure

Reduced set of privileges to manage the client. It is like the `manage` scope except the admin is not allowed to define protocol mappers, change the client template, or the client’s scope.

map-roles

Policies that decide if an admin can map any role defined by the client to a user. This is a shortcut, easy-of-use feature to avoid having to define policies for each and every role defined by the client.

map-roles-composite

Policies that decide if an admin can map any role defined by the client as a composite to another role. This is a shortcut, easy-of-use feature to avoid having to define policies for each and every role defined by the client.

map-roles-client-scope

Policies that decide if an admin can map any role defined by the client to the scope of another client. This is a shortcut, easy-of-use feature to avoid having to define policies for each and every role defined by the client.

##### [](#_users-permissions)Users

When going to the `Permissions` tab for all users, you will see these permission types listed.

view

Policies that decide if an admin can view all users in the realm.

manage

Policies that decide if an admin can manage all users in the realm. This permission grants the admin the privilege to perform user role mappings, but it does not specify which roles the admin is allowed to map. You’ll need to define the privilege for each role you want the admin to be able to map.

map-roles

This is a subset of the privileges granted by the `manage` scope. In this case the admin is only allowed to map roles. The admin is not allowed to perform any other user management operation. Also, like `manage`, the roles that the admin is allowed to apply must be specified per role or per set of roles if dealing with client roles.

manage-group-membership

Similar to `map-roles` except that it pertains to group membership: which groups a user can be added or removed from. These policies just grant the admin permission to manage group membership, not which groups the admin is allowed to manage membership for. You’ll have to specify policies for each group’s `manage-members` permission.

impersonate

Policies that decide if the admin is allowed to impersonate other users. These policies are applied to the administrator’s attributes and role mappings.

user-impersonated

Policies that decide which users can be impersonated. These policies will be applied to the user being impersonated. For example, you might want to define a policy that will forbid anybody from impersonating a user that has admin privileges.

##### [](#group)Group

When going to the `Permissions` tab for a specific group, you will see these permission types listed.

view

Policies that decide if the admin can view information about the group.

manage

Policies that decide if the admin can manage the configuration of the group.

view-members

Policies that decide if the admin can view the user details of members of the group.

manage-members

Policies that decide if the admin can manage the users that belong to this group.

manage-membership

Policies that decide if an admin can change the membership of the group. Add or remove members from the group.

## [](#_managing_organizations)Managing organizations

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/assembly-managing-organizations.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fassembly-managing-organizations.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fassembly-managing-organizations.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When integrating with a third party like a customer or business partner, you might want to manage their identities separately from others and build a unified and secure experience throughout your business ecosystem when they interact with a realm.

In a realm, an **organization** represents these third parties so that a realm or an organization administrator can manage the entire lifecycle of its members and how they authenticate and authorize to a realm, on a per-organization basis.

The organization is the entry point to start using the IAM capabilities of Keycloak to also address Business-to-Business (B2B) use cases. It enables multi-tenancy within a realm so that users can have access to protected resources from a realm but with a more restricted and controlled context, that context being the organization to which they belong.

Keycloak Organizations is a feature that enables support for organizations in Keycloak. For now, it provides some of the core capabilities needed to manage organizations such as:

- Manage members
- Organize members into groups with hierarchical structures
- Onboard organization members using invitation links
- Onboard organization members by federating their identities through identity brokering
- Identity-first login and organization-specific steps when authenticating in the scope of an organization
- Propagate organization-specific claims to applications through tokens for authorization purposes

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/organizations/intro.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Forganizations%2Fintro.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Forganizations%2Fintro.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

### [](#_enabling_organization_)Enabling organizations in Keycloak

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/organizations/managing-organization.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-organization.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-organization.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To use organizations, you have to enable the feature for the current realm.

Procedure

1. Click **Realm Settings** in the menu.
2. Toggle **Organizations** to **On**.
3. Click **Save**

Enabling Organizations

![Enabling Organizations](./images/organizations-enabling-orgs.png)

Once the feature is enabled, you are able to manage organizations through the **Organizations** section available from the menu.

### [](#managing-an-organization)Managing an organization

From the **Organizations** section, you can manage all the organizations in your realm.

Managing organizations

![Managing organizations](./images/organizations-management-screen.png)

#### [](#creating-an-organization)Creating an organization

Procedure

1. Click **Create Organization**.

Creating organization

![Creating organization](./images/organizations-create-org.png)

An organization has the following settings:

Name

A user-friendly name for the organization. The name is unique within a realm.

Alias

An alias for this organization, used to reference the organization internally. The alias is unique within a realm and must be URL-friendly, so characters not usually allowed in URLs will not be allowed in the alias. If not set, Keycloak will attempt to use the name as the alias. If the name is not URL-friendly, you will get an error and will be asked to specify an alias. Once defined, the alias cannot be changed afterwards.

Redirect URL

After completing registration or accepting an invitation to the organization sent via email, the user is automatically redirected to the specified redirect url. If left empty, the user will be redirected to the account console by default.

Domains

A set of zero or more domains that belongs to this organization. A domain cannot be shared by different organizations within a realm. When no domain is specified, organization members will not be validated against domain restrictions during authentication and profile validation.

Description

A free-text field to describe the organization.

Once you create an organization, you can manage the additional settings that are described in the following sections:

- [Manage attributes](#_managing_attributes_)
- [Manage members](#_managing_members_)
- [Manage groups](#_managing_groups_)
- [Manage identity providers](#_managing_identity_provider_)

#### [](#understanding-organization-domains)Understanding organization domains

When managing an organization, the domain associated with an organization plays an important role in how organization members authenticate to a realm and how their profiles are validated.

One of the key roles of a domain is to help to identify the organizations where a user is a member. By looking at their email address, Keycloak will match a corresponding organization using the same domain and eventually change the authentication flow based on the organization requirements.

The domain also allows organizations to enforce that users are not allowed to use a domain in their emails other than those associated with an organization. This restriction is especially useful when users, and their identities, are federated from identity providers associated with an organization and you want to force a specific email domain for their email addresses.

#### [](#disabling-an-organization)Disabling an organization

To disable an organization, toggle **Enabled** to **Off**.

Disabling organization

![Disabling organization](./images/organizations-disable-org.png)

When an organization is disabled, you can still manage it through the management interfaces, but the organization members cannot authenticate to the realm, including authenticating through the identity providers associated with the organization as they are also automatically disabled.

However, the unmanaged members of an organization are still able to authenticate to the realm as they are also realm users, but tokens will not hold metadata about their relationship with an organization that is disabled.

For more details about managed and unmanaged users, see [Managed and unmanaged members](#_managed_unmanaged_members_) section.

#### [](#deleting-an-organization)Deleting an organization

To delete an organization, click the **Delete** action for the corresponding organization in the listing page or when editing an organization.

Deleting organization

![Deleting organization](./images/organizations-delete-org.png)

When removing an organization, all data associated with it will be deleted, including any managed member.

Unmanaged users and identity providers remain in the realm, but they are no longer linked to the organization.

For more details about managed and unmanaged users, see [Managed and unmanaged members](#_managed_unmanaged_members_).

### [](#_managing_attributes_)Managing attributes

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/organizations/managing-attributes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-attributes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-attributes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

An administrator can store additional metadata about an organization using attributes. An organization attribute is a key/value pair that can hold multiple string values.

For that, click the **Attributes** tab and set any attribute you want by providing a key and a value.

To provide multiple values for the same attribute, and key, just add another attribute with the same key but with a different value.

Managing organization attributes

![Managing organization attributes](./images/organizations-manage-attributes.png)

### [](#_managing_members_)Managing members

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/organizations/managing-members.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-members.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-members.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

An organization member is basically a realm user but with a link to one or more organizations. They are logically separated from other users in a realm so that you know exactly which users belong to an organization.

There are different ways to add, or onboard, a member to an organization:

- Adding an existing realm user as a member
- Through an identity provider associated with an organization
- Sending an invitation to create a new account
- Sending an invitation to an existing user to join an organization

Once a member of an organization, that user’s account can be managed just like any regular account in a realm by accessing the **Users** section in the menu.

However, you can narrow the users to only those associated with an organization by accessing the **Members** tab when managing an organization. In this tab, you have a list of all the organization members and actions to add new members and to edit and remove existing ones.

Managing organization members

![Managing organization members](./images/organizations-manage-members.png)

You are also able to manage the membership of a user in an organization from the **Users** section. When you select a user, you can see the list of organizations that the user belongs to and manage those memberships.

Managing organization members

![Managing organization members from the Users section](./images/organizations-users-section.png)

#### [](#_managed_unmanaged_members_)Managed and unmanaged members

When managing members, consider how their relationship with an organization affects the lifecycle of their accounts. Members can join an organization through different flows and each flow indicates the strength of the link between their accounts and the organization.

There are two types of members:

- **Managed**
- **Unmanaged**

Managed members are those managed by the organization, and they cannot exist outside of their organization. For instance, consider an account created through an identity provider associated with an organization. That account does not belong to a realm as it was federated from the organization. In this case, the single source of truth for the identity is the organization and its lifecycle is controlled by the organization. If you remove the organization or the member, the account is also removed from the realm.

On the other hand, unmanaged members are those that can exist without the organization. For instance, when adding an existing realm user to an organization, the account belongs to the realm first and foremost and eventually linked to an organization. In this case, removing an organization or a member will not remove the account from the realm; the realm is the single-source of truth for the identity.

#### [](#adding-an-existing-realm-user-as-a-member)Adding an existing realm user as a member

An existing realm user can join an organization by selecting that user from a list and adding the user to the organization.

Procedure

1. Click **Add member**.
2. Click **Add realm user**.
3. Select one or more users and click **Add** to add them to the organization.

Adding a realm user

![Adding a realm user](./images/organizations-add-realm-user.png)

Once a user is a member of the organization, that user is able to authenticate to the realm just like a regular user and using any credential supported by the realm.

#### [](#inviting-users)Inviting users

An administrator can email users to join an organization.

Procedure

1. Click **Add member**.
2. Click **Invite member**.
3. Provide an email address
4. Click **Send**

Inviting members

![Inviting members](./images/organizations-invite-member.png)

Optionally, you can also provide a value for the **First name** and **Last name** fields for a more personalized email message using a greeting message with the first and last names of the person receiving the email.

An invitation is basically an email sent with a link that the person should click to perform the necessary steps to join an organization. These steps depend on whether the person already has an account in the realm or if a new account should be created before joining the organization.

If the email maps to an existing user in a realm, the steps the user will follow are basically about confirming the intention to join the organization.

On the other hand, if no user is associated with the given email address, the steps will involve creating a new account through the realm’s self-registration flow. In this case, the user is forced to provide the same email address used to send the invitation.

#### [](#managing-invitations)Managing invitations

Keycloak provides comprehensive invitation management capabilities that allow administrators to track, manage, and control organization invitations throughout their lifecycle.

##### [](#viewing-invitations)Viewing invitations

To view and manage all invitations for an organization:

Procedure

1. Select the organization you want to manage.
2. Click the **Invitations** tab.

The invitations list displays all sent invitations with the following information:

- **Email** - The email address of the invited user
- **Status** - Current status of the invitation (Pending, Expired)
- **Created** - When the invitation was sent
- **Expires** - When the invitation will expire (if applicable)

You can filter invitations by status to easily find pending or expired invitations.

##### [](#managing-invitation-lifecycle)Managing invitation lifecycle

Administrators can perform various actions on invitations:

###### [](#resending-invitations)Resending invitations

To resend a pending invitation:

Procedure

1. From the invitations list, locate the invitation you want to resend.
2. Click the action menu next to the invitation.
3. Select **Resend**.

A new invitation email will be sent to the recipient with a fresh expiration time.

###### [](#deleting-invitations)Deleting invitations

To permanently delete an invitation record:

Procedure

1. From the invitations list, locate the invitation you want to delete.
2. Click the action menu next to the invitation.
3. Select **Delete**.
4. Confirm the deletion.

This permanently removes the invitation from the system. This action cannot be undone.

##### [](#invitation-states-and-lifecycle)Invitation states and lifecycle

Invitations go through several states during their lifecycle:

- **Pending** - The invitation has been sent and is waiting for the recipient to accept
- **Expired** - The invitation has passed its expiration time

When a user successfully accepts an invitation, the invitation is automatically deleted from the system.

##### [](#api-access)API access

Organization invitations can also be managed programmatically through the Admin REST API:

- `GET /admin/realms/{realm}/orgs/{orgId}/invitations` - List all invitations
- `GET /admin/realms/{realm}/orgs/{orgId}/invitations/{invitationId}` - Get specific invitation
- `POST /admin/realms/{realm}/orgs/{orgId}/invitations/{invitationId}/resend` - Resend invitation
- `DELETE /admin/realms/{realm}/orgs/{orgId}/invitations/{invitationId}` - Delete invitation

For detailed API documentation, refer to the Keycloak Admin REST API documentation.

#### [](#_onboard_member_identity_provider_)Onboarding members using an Identity Provider

An organization might have its own identity provider as the single source of truth for their identities. In this case, users federated from the identity provider are automatically added as a member of the organization.

When users join an organization through an identity provider associated with an organization, they are automatically marked as managed members. In this case, they will go through the broker login flows configured in the realm and join the organization automatically once they successfully authenticate.

Onboarding new members through an identity provider can be done by either automatically redirecting the user to an organization’s identity provider or by selecting the identity provider when at the login page.

In both cases, once the user provides the email, Keycloak will try to match an organization based on the email domain. In case the email domain matches the organization, and an identity provider is associated with the same domain and the **Redirect when email domain matches** setting is enabled, the user is automatically redirected to the identity provider. Once the user authenticates at the identity provider and completes the first broker login flow, the user is automatically added as an organization member.

On the other hand, if **Redirect when email domain matches** is not enabled, but the identity provider is configured not to **Hide on login page**, the user can select the identity provider and then be redirected to the identity provider to continue the onboarding process.

For more details, see [Managing Identity Providers](#_managing_identity_provider_).

#### [](#removing-a-member)Removing a member

You can remove a member from an organization.

From the action menu next to the member you want to remove, click **Remove**.

When removing a member from an organization, remember that the user may or may not be removed from a realm depending on if that user is managed or unmanaged member, respectively.

For more details, see [Managed and unmanaged members](#_managed_unmanaged_members_).

#### [](#viewing-organization-group-memberships)Viewing organization group memberships

You can view which organization groups a member belongs to directly from the member’s context menu.

Procedure

1. To view group memberships from the **Members** tab:
   
   1. Navigate to your organization.
   2. Click the **Members** tab.
   3. Locate the member.
   4. Click the kebab menu (⋮) next to the member.
   5. Select **Show organization group memberships**.
2. To view group memberships from a specific group:
   
   1. Navigate to your organization.
   2. Click the **Groups** tab.
   3. Select the group.
   4. Click the **Members** tab.
   5. Locate the member.
   6. Click the kebab menu (⋮) next to the member.
   7. Select **Show memberships**.
3. To view group memberships from the user’s detail page:
   
   1. Click **Users** in the menu.
   2. Select the user.
   3. Click the **Organizations** tab.
   4. Locate the organization.
   5. Click the kebab menu (⋮) next to the organization.
   6. Select **Show organization group memberships**.

This displays all groups from the specific organization the member is currently assigned to.

#### [](#support-for-federated-members)Support for federated members

Users coming from federated providers can also be added as members of an organization. The only exceptions are the users from LDAP providers with **import mode disabled**. Organization members are added to an internal group that is not synchronized with external providers, so even if the LDAP provider has a group mapper with mode LDAP\_ONLY it won’t be possible for the non-imported users to be added as members of an organization because that membership won’t be synced with the LDAP server.

In other words, LDAP users that are not imported can’t join an organization because the membership is not stored in the local DB nor in the LDAP server. So if you want to have LDAP users joining organizations, ensure that the import mode of the LDAP provider is enabled.

### [](#_managing_groups_)Managing groups

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/organizations/managing-groups.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-groups.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-groups.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Organization groups let you organize members into logical teams, departments, or any hierarchical structure that makes sense for your organization. Think of them as folders for your users—simple, flexible, and powerful.

Unlike realm groups that are shared across your entire realm, organization groups belong exclusively to one organization. This means `orgA` and `orgB` can each have a `/Engineering/Backend` group without any conflict—each is a completely separate group with its own unique identifier.

#### [](#creating-groups)Creating groups

Navigate to your organization and click the **Groups** tab.

Managing organization groups

![Managing organization groups](./images/organizations-manage-groups.png)

Procedure

1. Click **Create group**.
2. Provide a **Name** for the group.
3. Click **Create**.
4. Optionally, before creating the group, navigate to **Parent group** to create nested structures.

Your group is ready. You can now add members, nest more groups under it, or set attributes.

##### [](#building-hierarchies)Building hierarchies

Groups can be nested to create organizational structures that mirror your real-world teams:

```
/Engineering
  /Engineering/Backend
  /Engineering/Frontend
/Sales
  /Sales/APAC
  /Sales/EMEA
  /Sales/LATAM
  /Sales/NA
```

Create a child group by selecting its parent before creation, or move groups, by clicking **Move to** in the UI to reorganize them.

#### [](#adding-members-to-groups)Adding members to groups

Members join groups the same way they join realm groups.

Procedure

1. Select the group.
2. Click the **Members** tab.
3. Click **Add member**.
4. Select one or more organization members.
5. Click **Add**.

To automatically assign members to these groups when they authenticate via an Identity Provider, see [Mapping federated users to organization groups](#_managing_identity_provider_).

To view which organization groups a specific member belongs to, use the kebab menu (⋮) next to the member and select **Show group memberships**. For details, see [Viewing organization group memberships](#viewing-organization-group-memberships).

#### [](#understanding-group-paths)Understanding group paths

Every group has a path that describes its location in the hierarchy. Paths are **relative to the organization**, not the realm.

For example:

- Organization A has `/Engineering/Backend`
- Organization B has `/Engineering/Backend`
- The realm also has `/Engineering/Backend`

These are three completely separate groups. When you reference `/Engineering/Backend` in the context of Organization A, you’re talking about Organization A’s group—period.

This isolation is intentional. Your organizations remain independent, and you can structure each one however you need without worrying about naming conflicts.

#### [](#mapping-groups-to-tokens)Mapping groups to tokens

Organization group membership can appear in tokens for both OIDC and SAML protocols. The setup differs between protocols.

##### [](#oidc-organization-groups-in-tokens)OIDC: Organization groups in tokens

Organization groups are added to the `organization` claim alongside other organization data. To include groups in OIDC tokens, you must configure the mapper AND request the appropriate scope.

The **Organization Group Membership** mapper alone is NOT sufficient. It only works when combined with the **Organization Membership** mapper in the same scope.

###### [](#option-1-add-mapper-to-built-in-organization-scope-recommended)Option 1: Add mapper to built-in organization scope (Recommended)

The built-in `organization` scope already has the Organization Membership mapper. Simply add the group mapper to it.

Procedure

1. Click **Client scopes** in the menu.
2. Select the **organization** scope.
3. Click the **Mappers** tab.
4. Click **Add mapper** → **By configuration**.
5. Select **Organization Group Membership**.
6. Configure the mapper:
   
   Add to ID token
   
   Include groups in ID tokens.
   
   Add to access token
   
   Include groups in access tokens.
   
   Add to lightweight access token
   
   Include groups in lightweight access tokens.
   
   Add to userinfo
   
   Include groups in the UserInfo endpoint response.
   
   Add to introspection
   
   Include groups in the Introspection endpoint response.
7. Click **Save**.

The `organization` scope is added as an optional scope to all clients by default. Applications request it using `scope=organization` (or `scope=organization:alias` for a specific organization).

###### [](#option-2-create-custom-scope-with-both-mappers)Option 2: Create custom scope with both mappers

If you need a dedicated scope for organization data, create one with BOTH required mappers.

Procedure

01. Click **Client scopes** in the menu.
02. Click **Create client scope**.
03. Provide a **Name** (for example, `my-org-scope`).
04. Set **Protocol** to **OpenID Connect**.
05. Click **Save**.
06. Click the **Mappers** tab.
07. Add BOTH mappers:
    
    1. Click **Configure a new mapper** → **Organization Membership**
    2. Configure and save
    3. Click **Add mapper** → **By configuration** → **Organization Group Membership**
    4. Configure and save
08. Go to your client.
09. Click **Client scopes** tab.
10. Click **Add client scope**.
11. Select your custom scope.
12. Click **Add**.
13. Choose **Default** or **Optional**.

Applications request it using `scope=my-org-scope`.

###### [](#token-structure)Token structure

When a user from Organization A’s `/Engineering/Backend` group authenticates and requests the organization scope, the token includes:

```
{
  "organization": {
    "orgA": {
      "id": "f8d3c4e1-...",
      "groups": [ "/Engineering/Backend" ]
    }
  }
}
```

Notice:

- Groups appear within the organization claim, not as a separate top-level claim
- Group paths are relative to the organization
- Multiple organizations can be included if the user is a member of multiple organizations and uses `scope=organization:*`

##### [](#saml-organization-groups-in-assertions)SAML: Organization groups in assertions

SAML automatically includes groups for all organizations the user is a member of. Like OIDC, the group mapper requires the organization membership mapper to be present in the same scope.

###### [](#option-1-add-mapper-to-built-in-saml_organization-scope-recommended)Option 1: Add mapper to built-in saml\_organization scope (Recommended)

The built-in `saml_organization` scope already has the Organization Membership mapper and is added as a default scope to all SAML clients. Simply add the group mapper to it.

Procedure

1. Click **Client scopes** in the menu.
2. Select the **saml\_organization** scope.
3. Click the **Mappers** tab.
4. Click **Add mapper** → **By configuration**.
5. Select **Organization Group Membership**.
6. Click **Save**.

The mapper has no additional configuration options. Groups are automatically included for all SAML clients.

###### [](#option-2-create-custom-saml-scope-with-both-mappers)Option 2: Create custom SAML scope with both mappers

If you need a dedicated scope, create one with both the membership and group mappers.

Procedure

01. Click **Client scopes** in the menu.
02. Click **Create client scope**.
03. Provide a **Name** (for example, `my-saml-org-scope`).
04. Set **Protocol** to **SAML**.
05. Click **Save**.
06. Click the **Mappers** tab.
07. Add BOTH mappers:
    
    1. Click **Configure a new mapper** → **Organization Membership**
    2. Click **Save**
    3. Click **Add mapper** → **By configuration** → **Organization Group Membership**
    4. Click **Save**
08. Go to your SAML client.
09. Click **Client scopes** tab.
10. Click **Add client scope**.
11. Select your custom scope.
12. Click **Add**.
13. Choose **Default** or **Optional**.

###### [](#saml-assertion-structure)SAML assertion structure

When a user from Organization A’s `/Engineering/Backend` group authenticates, the assertion includes attributes for each organization they belong to:

```
<Attribute Name="organization.orgA.groups">
  <AttributeValue>/Engineering/Backend</AttributeValue>
</Attribute>
```

Each organization gets its own attribute named `organization.{alias}.groups` with group paths as values, where `alias` is the organization alias.

Unlike OIDC, SAML clients don’t need to request a scope at runtime. The `saml_organization` scope is added as a default scope, so groups are automatically included in assertions for all organizations the user is a member of.

#### [](#managing-group-attributes)Managing group attributes

Groups can carry attributes—key/value metadata stored at the group level.

Procedure

1. Select a group.
2. Click the **Attributes** tab.
3. Add attributes with keys and values.
4. Click **Save**.

Group attributes are stored separately from user attributes. To include group attributes in tokens, configure a **User Attribute** mapper with the **Aggregate attributes** option enabled. The mapper will then combine matching attributes from both the user and all their groups.

#### [](#important-distinctions)Important distinctions

##### [](#organization-groups-vs-realm-groups)Organization groups vs realm groups

  Feature Organization groups

Scope

Belong to a single organization

Isolation

Fully isolated—same paths can exist across organizations

Role assignment

Not supported (coming in future releases)

Authorization policies

Cannot be used in authorization policies—use realm groups instead

Token mapping

Appear in tokens when using organization context

##### [](#authorization-policies-restriction)Authorization policies restriction

Organization groups **cannot** be used in Keycloak authorization policies. Only realm groups are supported for fine-grained authorization.

If you try to select an organization group when configuring a group-based policy, you’ll receive an error. This is intentional—authorization policies operate at the realm level, not the organization level.

#### [](#deleting-groups)Deleting groups

Deleting a group removes it and all its child groups. Members are **not** deleted—they simply lose their group membership.

Procedure

1. Select the group.
2. Click **Delete group**.
3. Confirm the deletion.

### [](#_managing_identity_provider_)Managing identity providers

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/organizations/managing-identity-providers.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-identity-providers.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Forganizations%2Fmanaging-identity-providers.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

An organization might have its own identity provider as the single source of truth for their identities. In this case, you want to configure the organization to authenticate users using the organization’s identity provider, federate their identities, and finally add them as a member of the organization.

An organization can have one or more identity providers associated with it so that they can authenticate their users from different sources and enforce different constraints on each of them.

Before you can link an identity provider to an organization, you create an organization at the realm level from the **Identity Providers** section in the menu. You can link any of the built-in social and identity providers available in the realm to an organization.

#### [](#linking-an-identity-provider-to-an-organization)Linking an identity provider to an organization

An identity provider can be linked to an organization from the **Identity providers** tab. If identity providers already exist, you see a list of them and options to search, edit, or unlink from the organization.

Organization identity providers

![Organization identity providers](./images/organizations-identity-providers.png)

Procedure

1. Click **Link identity provider**
2. Select an **Identity provider**
3. Set the appropriate settings
4. Click **Save**

Linking identity provider

![Linking identity provider](./images/organizations-link-identity-provider.png)

An identity provider has the following settings:

Identity provider

The identity provider you want to link to the organization. An identity provider can only be linked to a single organization.

Domain

The domain from the organization that you want to link with the identity provider.

Hide on login page

If this identity provider should be hidden in login pages when the user is authenticating in the scope of the organization.

Hide on login page when organization not resolved

If enabled, the identity provider will be hidden on the login page when the organization cannot be resolved based on the user’s email domain. Otherwise, the identity provider will be shown on the login page regardless of whether the organization is resolved or not. If 'Hide on login page' is also enabled, the identity provider will always be hidden on the login page.

Redirect when email domain matches

If members should be automatically redirected to the identity provider when their email domain matches the domain set to the identity provider. If the domain is set to `Any`, members whose email domain matches **any** of the organization domains will be redirected to the identity provider.

If the org is linked with multiple identity providers, the organization authenticator prioritizes the provider that matches the email domain of the user for automatic redirection. If none is found, it tries to locate one whose domain is set to `Any`.

Once linked to an organization, the identity provider can be managed just like any other in a realm by accessing the **Identity Providers** section in the menu. However, the options herein described are only available when managing the identity provider in the scope of an organization. The only exception is the **Hide on login page** option that is present here for convenience.

##### [](#mapping-federated-users-to-organization-groups)Mapping federated users to organization groups

When an identity provider is linked to an organization, you can configure mappers to automatically assign users authenticating through the IdP to organization groups based on claims or attributes from the external IdP.

###### [](#using-the-hardcoded-group-mapper)Using the Hardcoded Group mapper

The Hardcoded Group mapper automatically adds all users authenticating through the IdP to a specific group.

Procedure

1. Click **Identity Providers** in the menu.
2. Select your identity provider linked to your organization.
3. Click the **Mappers** tab.
4. Click **Add mapper**
5. Select **Hardcoded Group** mapper type.
6. Configure the mapper:
   
   Name
   
   A descriptive name for this mapper.
   
   Sync Mode Override
   
   How membership is managed (Import, Force, Inherit).
   
   Group
   
   The group to assign users to. If the IdP is linked to an organization, you can either select realm groups or groups from that organization.
7. Click **Save**.

Groups are filtered to the organization linked with this identity provider. If the IdP is unlinked from the organization, users will no longer be assigned to organization groups during authentication.

###### [](#using-the-advanced-group-mapper)Using the Advanced Group mapper

The Advanced Group mapper assigns users to groups based on claim values from the external IdP.

Procedure

1. Click **Identity Providers** in the menu.
2. Select your identity provider linked to your organization.
3. Click the **Mappers** tab.
4. Click **Add mapper**
5. Select **Advanced Claim to Group**.
6. Configure the mapper:
   
   Name
   
   A descriptive name for this mapper.
   
   Sync Mode Override
   
   How membership is managed (Import, Force, Inherit).
   
   Claims
   
   Key-value pairs to match claims from the external IdP. Add claim names and their expected values. Users matching these claim values will be assigned to the specified group.
   
   Regex Claim Values
   
   When enabled, claim values are treated as regular expressions for pattern matching. When disabled, values must match exactly.
   
   Group
   
   The target group. If the IdP is linked to an organization, you can select groups from that organization alongside with realm groups.
7. Click **Save**.

###### [](#group-selection-behavior)Group selection behavior

When configuring group mappers for an identity provider:

- **IdP linked to organization:** Both realm groups and groups from that organization are selectable
- **IdP not linked to organization:** Only realm groups are selectable
- **IdP unlinked from organization:** Existing organization group mappings are preserved but no longer active.

###### [](#runtime-validation)Runtime validation

When a user authenticates through an identity provider:

1. The mapper checks if the IdP is still linked to an organization
2. If the IdP is unlinked, organization group mappings are skipped
3. Only realm group mappings continue to work
4. A warning is logged to help administrators identify outdated mapper configurations
5. Authentication continues successfully—no error is thrown to the user

This ensures that unlinking an IdP from an organization doesn’t break authentication, but organization-specific mappings become inactive.

###### [](#migration-considerations)Migration considerations

If you unlink an identity provider from an organization:

- Review all mappers on that IdP
- Update any mappers referencing organization groups
- Consider changing those mappings to realm groups if needed
- Alternatively, re-link the IdP to the organization to restore group mappings

Organization groups cannot be shared across organizations. If you move an IdP from one organization to another, update the group mappings accordingly.

#### [](#editing-a-linked-identity-provider)Editing a linked identity provider

You can edit any of the organization-related settings of a linked identity provider at any time.

Procedure

1. In the menu, click **Organizations** and go to the **Identity providers** tab.
2. Locate the **identity provider** in the list.
   
   You can use the search option for this step.
3. Click the action button (three dots) at the end of the line.
4. Click **Edit**.
5. Make the necessary changes.
6. Click **Save**.

Editing linked identity provider

![Editing linked identity provider](./images/organizations-edit-identity-provider.png)

#### [](#unlinking-an-identity-provider-from-an-organization)Unlinking an identity provider from an organization

When an identity provider is unlinked from an organization, it remains available as a realm-level provider that is no longer associated with an organization. To delete the unlinked provider, use the **Identity Providers** section in the menu.

Procedure

1. In the menu, click **Organizations** and go to the **Identity providers** tab.
2. Locate the **identity provider** in the list.
   
   You can use the search capabilities for this step.
3. Click the action button (three dots) at the end of the line.
4. Click **Unlink provider**.

Unlinking identity provider

![Unlinking identity provider](./images/organizations-unlink-identity-provider.png)

### [](#authenticating-members_server_administration_guide)Authenticating members

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/organizations/authenticating-members.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Forganizations%2Fauthenticating-members.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Forganizations%2Fauthenticating-members.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When you enable organizations for a realm, user authentication is changed. If the user is recognized to be authenticating in the context of an organization, the authentication flow changes on a per-organization basis.

When a realm is created, the authentication flows are automatically updated to enable specific steps to authenticate and onboard organization members. The authentication flows updated are:

- **browser**
- **first broker login**

The main change to the **browser** flow is that it defaults to an identity-first login so that users are identified before prompting for their credentials. Concerning the first broker login flow, the main change is automatically adding the users as organization members once they authenticate through the identity provider associated with an organization and successfully complete the flow.

The choice to use an identity-first login or not depends on the existence of an organization in a realm. If no organizations exist, the user follows the usual steps to authenticate using the username and password, or any other step configured in the browser flow. Otherwise, the user is asked first for a username or email to continue authenticating to a realm.

The identity-first login main goal is to identify the user:

- Is the user an existing or a new user?
- Is the user a member of any organization within a realm?
- If an organization member, is the user linked to any identity provider associated with the organization?

Depending on the outcome when identifying the user, the authentication flow changes to either proceed with authentication by asking for the user’s credentials or eventually redirect the user automatically to authenticate within the organization security boundaries through an identity provider.

#### [](#understanding-the-identity-first-login)Understanding the identity-first login

In addition to identifying the user once the username is provided, the identity-first login is also responsible for:

- Matching an email domain to an organization.
- Deciding if the authentication flow should continue or not if an account already exists for the username provided
- Deciding how the user should be authenticated depending on how the domains and the identity providers are configured to an organization and the set of credentials configured to the user account.
- Seamlessly authenticating users through an identity provider associated with an organization if the email domain matches the domain set to the identity provider

The identity-first login provides the same capabilities that are provided by the usual login page with the username and password fields. Users can still self-register by clicking the register link or choose any identity or social broker that is not linked to an organization in that realm.

Identity-first login page

![Identity-first login page](./images/organizations-identity-first-login.png)

In the case of a user that does not exist, if that user tries to authenticate using an email domain that matches an organization domain, the identity-first login page appears again with a message that the username provided is not valid. At this point, no need exists to ask the user for credentials.

Identity-first when user does not exist

![Identity-first login error](./images/organizations-identity-first-error.png)

Several options exist to register the user allowing that user to authenticate to the realm and join an organization.

If the realm has the self-registration setting enabled, the user can click the **Register** link at the identity-first login page and create an account at the realm. After that, the administrator can send an invitation link to the user or manually add the user as a member of an organization. For more details, see [Managing members](#_managing_members_).

If the organization has an identity provider without a domain and the **Hide on login page** setting is **OFF**, users can also click the identity provider link at the identity-first login page to automatically create an account and join an organization once they authenticate through the identity provider. For more details, see [Managing identity providers](#_managing_identity_provider_).

In a similar situation to the previous section, the organization may have an identity provider set with one of the organization domains. In this situation, the user is redirected to the identity provider if that user’s email matches a specific domain from the organization. Once the flow completes, an account is created and the user joins the organization. In case the user has any first-factor credentials configured (e.g.: password, passwordless, kerberos) to the account, the user is not automatically redirected to the identity provider but asked to authenticate using their credentials.

#### [](#configuring-existing-authentication-flows)Configuring existing authentication flows

As previously mentioned for new realms, authentication flows are automatically updated with the necessary steps to support organizations and authenticate their members. For existing realms, in addition to enabling organizations to the realm, you also need to manually update your existing (custom) authenticating flows.

Change the **browser** flow by following these steps:

Procedure

1. Duplicate the current flow bound to the **Browser flow** binding type to avoid breaking the flow you are currently using
2. Click **Add sub-flow** and give it a name such as **My Organization**
3. Move the newly added **My Organization** sub-flow to execute right after the **Identity Provider Redirector** execution step. The main point here is that the sub-flow should happen before any other sub-flow or execution step that authenticates the user using whatever credentials you support in your realm. Once added, change the **Requirement** to **Alternative**.
4. Click **Add sub-flow** in the **My Organization** sub-flow and give it a name such as **My Organization - Conditional**. Once added, change the **Requirement** to **Conditional**.
5. Click **Add condition** in the **My Organization - Conditional** sub-flow and select **Condition - user configured**. Once added, change the **Requirement** to **Required**.
6. Click **Add step** in the **My Organization - Conditional** sub-flow and select the \*Organization Identity-First Login
   
   - execution step. Once added, change the **Requirement** to **Alternative**.
7. Bind the authentication flow to the **Browser** binding type.

Organizations browser flow

![Organizations browser flow](./images/organizations-browser-flow.png)

Once you enable the [Organizations](#_enabling_organization_) setting to the realm and create at least a single organization, you should be able to see the identity-first login page and start using organizations in your realm.

Change the **first broker login** flow by following these steps:

Procedure

1. Duplicate the current flow bound to the **First broker login flow** bind type to avoid breaking the flow you are currently using
2. Click **Add sub-flow** and give it a name such as `Organization Member - Conditional`. Once added, change the **Requirement** to **Conditional**.
3. Click **Add condition** in the **Organization Member - Conditional** sub-flow and select **Condition - user configured**. Once added, change the **Requirement** to **Required**.
4. Click **Add step** in the **Organization Member - Conditional** sub-flow and select the **Organization Member Onboard** execution step. Once added, change the **Requirement** to **Required**.
5. Bind the authentication flow to the **First broker login** binding type.

Organizations first broker flow

![Organizations first broker flow](./images/organizations-first-broker-flow.png)

You should now be able to authenticate using any identity provider associated with an organization and have the user joining the organization as a member as soon as they complete the first browser login flow.

#### [](#configuring-how-users-authenticate)Configuring how users authenticate

If the flow supports organizations, you can configure some of the steps to change how users authenticate to the realm.

For example, some use cases will require users to authenticate to a realm only if they are a member of any or a specific organization in the realm.

To enable this behavior, you need to enable the `Requires user membership` setting on the `Organization Identity-First Login` execution step by clicking on its settings.

If enabled, and after the user provides the username or email in the identity-first login page, the server will try to resolve a organization where the user is a member by looking at any existing membership or based on the semantics of the [organization](#_mapping_organization_claims_) scope, if requested by the client. If not a member of an organization, an error page will be shown.

### [](#_mapping_organization_claims_)Mapping organization claims

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/organizations/mapping-organization-claims.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Forganizations%2Fmapping-organization-claims.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Forganizations%2Fmapping-organization-claims.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To map organization-specific claims into tokens, a client needs to request the **organization** scope when sending authorization requests to the server. When authenticating in the context of an organization, clients can request the `organization` scope to map information about the organizations where the user is a member.

As a result, the token will contain a claim as follows:

```
"organization": {
  "testcorp": {
    "id": "42c3e46f-2477-44d7-a85b-d3b43f6b31fa",
    "attr1": [
      "value1"
    ]
  }
}
```

The organization claim can be used by clients (for example, from ID Tokens) and resource servers (for example, from access tokens) to authorize access to protected resources based on the organization where the user is a member.

The `organization` scope is a built-in optional client scope at the realm. Therefore, this scope is added to any client created in the realm by default. It also defines the `Organization Membership` mapper that controls how the organization membership information is mapped to the tokens.

By default, the organization id and attributes are not included in the organization claim. To include them, edit the mapper and enable the **Add organization id** and **Add organization attributes** options, respectively.

Including attributes in the organization claim

![Including attributes in the organization claim](./images/organizations-add-org-attrs-in-claim.png)

The `organization` scope is requested using different formats:

  Format Description

`organization`

Maps to a single organization if the user is a member of a single organization. Otherwise, if a member of multiple organizations, the user will be prompted to select an organization when authenticating to the realm.

`organization:<alias>`

Maps to a single organization with the given alias.

`organization:*`

Maps to all organizations the user is a member of.

For details on including organization groups in tokens, see [Managing groups](#_managing_groups_).

## [](#_managing_workflows)Managing workflows

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/assembly-managing-workflows.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fassembly-managing-workflows.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fassembly-managing-workflows.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak Workflows is a powerful engine to automate and orchestrate realm administrative tasks, bringing key capabilities of Identity Governance and Administration (IGA) to your identity and access management infrastructure. By using workflows, you can implement policies and processes that govern the lifecycle of realm resources, such as users and clients, helping you to improve security, meet compliance requirements, and reduce administrative costs.

As a core component of IGA, identity lifecycle management is fully supported by workflows so that you can easily automate onboarding and offboarding processes, and other recurring administrative tasks. For example, you can define workflows to:

- Provision and de-provision realm resources, such as users, automatically based on specific events or conditions.
- Automate Joiner-Mover-Leaver (JML) processes to ensure that users are granted the appropriate access rights based on their roles and responsibilities.
- Enforce access reviews and certifications to ensure that users have the appropriate access rights.
- Enforce Just-In-Time and Least Privilege access by automating role assignments and revocations.
- Enforce strong authentication policies for realm users.
- Prevent insider data breaches by automating the removal of inactive or obsolete users.

By leveraging workflows, realm administrators can ensure that security policies are consistently enforced and based on the least privilege principle, reduce the risk of human error, and free up valuable time to focus on other important administrative tasks. This guide will walk you through the process of creating and managing workflows to automate your administrative tasks and implement IGA best practices when managing realms.

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/intro.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Fintro.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Fintro.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

### [](#_understanding_workflows_)Understanding workflows

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/understanding-workflow.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-workflow.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-workflow.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

A workflow is an activity or process within Keycloak that executes a series of predefined steps on a specific realm resource in response to specific events and based on the defined conditions.

In Keycloak, there are two main types of events:

- **User Events**: These events are initiated by users, such as login in or registering a new account.
- **Admin Events**: These events are initiated by administrators through the Admin Console and API, such as creating or updating a user or any other realm resource.

Workflows are designed to respond to these events by automating tasks and processes that need to be performed and executed on the realm resources associated with these events.

Once an event is fired, the workflow engine evaluates all defined workflows in the realm to determine if any should be triggered based on the event type and the specified conditions. If the event matches the conditions defined in a workflow, a workflow execution is created and assigned with an identifier. The workflow execution is bound to the realm resource associated with the event, and only a single instance of a workflow can be active for that realm resource at any given time.

A realm resource can be any entity within the realm, such as a user, client, group, or a role. At the moment, workflows can be defined for the following realm resources:

- Users

Once a workflow is active for a realm resource, its step chain is processed. The steps run sequentially, and they will run immediately or be scheduled to execute later in time, depending on each step configuration. Once all the steps are executed, the workflow execution is completed and the realm resource is no longer bound to that workflow.

That is the main gist of Keycloak Workflows. There are more details and settings that can be configured to customize the behavior and the lifecycle of a workflow, how their instances are created, managed, and bound to realm resources, as well as how their steps are executed.

### [](#_understanding_workflow_definition_)Understanding the workflow definition

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/understanding-workflow-definition.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-workflow-definition.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-workflow-definition.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Workflows are defined in YAML format. This format allows for a clear and human-readable way to specify the automation process represented by a workflow.

In its simplest form, a workflow definition can be defined as follows:

```
name: Onboarding new users
on: user-created
steps:
  - uses: notify-user
    with:
      message: |
        <p>Dear ${user.firstName} ${user.lastName}, </p>

        <p>Welcome to ${realm.displayName}!</p>

        <p>
           Best regards,<br/>
           ${realm.displayName} Team
        </p>
  - uses: add-required-action
    after: 30d
    with:
      action: UPDATE_PASSWORD
  - uses: restart
    with:
      position: 1
```

The example above is very simple but illustrates the core structure and the flexibility of a workflow to automate administrative tasks. It is composed of three main sections:

- `name`: A unique and user-friendly name to identify the workflow within the realm.
- `on`: The event that will trigger the workflow. In this case, the workflow is triggered when a new user is added to the realm.
- `steps`: A set of one or more steps to be executed when executing a workflow execution. In this example, three steps are defined:
  
  1. The first step uses the built-in `notify-user` action to send a welcome message to the new user.
  2. The second step uses the built-in `add-required-action` action to require the user to update their password after 30 days.
  3. The third step uses the built-in `restart` action to restart the workflow from the second step so that the user is forced to update their password every 30 days.

Here is a more detailed look at all settings available from the workflow definition:

`name`

A unique and user-friendly name to identify the workflow within the realm. This name is crucial for management and logging purposes. This setting is mandatory.

`on`

Define a condition that determines the event that will trigger the workflow. The condition is written using an expression language that supports a variety of checks on the event. See [Triggering workflows on events](#_workflow_events_) for more details. This setting is optional.

`schedule`

Define a schedule that will trigger the workflow at defined intervals. See [Scheduling workflows](#_scheduling_workflows_) for more details. This setting is optional.

`if`

Define a condition that must be met for the workflow to be triggered. The condition is written using an expression language that supports a variety of checks on the realm resource associated with the event. A workflow execution is only created if the expression evaluates to `true`. If this setting is omitted, the event defined in the `on` setting will always create the workflow execution. See [Defining conditions](#_workflow_conditions_) for more details. This setting is optional.

`steps`

Define the step chain consisting of one or more steps to be sequentially executed during the lifetime of a workflow. See [Defining steps](#_workflow_steps_) for more details. This setting is mandatory.

`concurrency`

This setting controls the behavior when multiple events that would trigger the same workflow occur for the same realm resource while a workflow execution is already active for that resource. This setting is optional. + The available options are:

- `restart-in-progress`: A new event causes the current workflow execution to restart from the beginning.
- `cancel-in-progress`: A new event causes the current workflow execution to be canceled and disassociated from the realm resource.

`enabled`

This setting enables or disables the workflow. If set to `false`, new workflow executions will not be created for the workflow and any active workflow execution will be paused. This setting is optional and defaults to `true`.

### [](#_workflow_expression_language_)Understanding the workflow expression language

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/understanding-workflow-expression.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-workflow-expression.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-workflow-expression.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Some settings from the workflow definition can be defined using a boolean expression language, the Workflow Expression Language. Expressions are defined using functions, operands, and logical operators. The functions available in an expression will depend on the setting where it is being defined.

The functions can have zero or more arguments, and they can be combined with the following logical operators and symbols to define complex expressions:

  Operator Description

`and`

Logical `AND`

`or`

Logical `OR`

`not`

Logical `NOT`

`()`

Grouping and delimiting conditional logic

These examples are all valid expressions:

```
f1
f1 and f2
not (f1 and f2)
(f1 and f2) or f3
```

The settings that supports expressions provide its own set of functions as you will see in the following sections.

### [](#_managing_workflows_)Managing workflows

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/managing-workflows.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Fmanaging-workflows.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Fmanaging-workflows.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Workflows can be managed through the Admin Console or the Admin REST API.

Only realm administrators with the appropriate permissions can manage workflows as they are considered sensitive operations. For more details, see [Understanding different types of Realm Admins](#_understanding_different_types_realm_admins_).

#### [](#managing-workflows-through-the-admin-console)Managing workflows through the Admin Console

To manage workflows through the Admin Console, you can follow these steps:

Procedure

1. Click **Workflows** in the menu.

Empty workflows list

![Empty workflows list](./images/workflows-empty-list.png)

If your realm does not have any workflows defined yet, you will see an empty list. To create a new workflow, click **Create**.

1. Creating a new workflow ![Creating a new workflow](./images/workflows-create.png)

This will open the workflow creation screen where you can define the workflow using YAML format. To save the workflow, click **Save**.

Once you have created one or more workflows, they will be listed in the workflows list.

List of workflows

![List of workflows](./images/workflows-list.png)

From the workflows list, you can perform the following actions:

- **Create**: Click **Create workflow** to create a new workflow.
- **Update**: Click on the name of an existing workflow to update it.
- **Enable/Disable**: Use the toggle button on the **Status** column to enable or disable a workflow.
- **Copy**: Click on the **Copy** button next to a workflow to create a copy of it. This is useful if you want to create a new workflow based on an existing one.
- **Delete**: Click on the **Delete** button next to a workflow to delete it.

##### [](#updating-a-workflow)Updating a workflow

To update a workflow, click on its name in the workflows list. This will open the workflow in edit mode, where you can modify its YAML definition. After making your changes, click **Save** to apply the updates.

There are limitations when updating workflows if they have scheduled steps. In such cases, you are not able to change the step chain of a workflow if there are active instances of the workflow being updated. This limitation is in place to ensure the integrity of ongoing workflow instances where updates to the step chain could lead to an expensive operation of migrating existing instances to the new step chain, which could impact system performance as well as lead to inconsistencies and an unexpected behavior.

In future releases, we plan to introduce a more flexible approach to updating workflows with scheduled steps, so that administrators are going to be allowed to migrate realm resources scheduled for a specific step to a different step in the updated workflow.

For now, the only option is to either wait for all active instances to complete or to delete the workflow and create a new one. Note that deleting a workflow will also delete all its active instances, therefore, realm resources associated with those instances will not be processed further, and they would need to be re-associated with a new workflow instance if needed.

### [](#_workflow_events_)Triggering workflows on events

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/listening-workflow-events.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Flistening-workflow-events.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Flistening-workflow-events.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The `on` setting defines the event that will trigger the workflow. You should choose the event accordingly to the realm resource and the intent of the workflow.

The value of the \`on' setting is an expression that must evaluate to true. In its basic form, the expression is simply the name of the event.

```
on: user-created
```

As you can see, events are functions in an expression and, as such, they also support parameters to provide more specific conditions. For example, to trigger a workflow when a user is added to a specific group, you can use the following expression:

```
on: user-group-membership-added(/mygroup)
```

You can also combine multiple conditions using logical operators. For example, to trigger a workflow when a user is created or when they authenticate, you can use:

```
on: user-created or user-authenticated
```

Even though the `on` setting can be defined using expressions, keep in mind that a single event triggers a workflow.

If this setting is not specified, the workflow will not be triggered by any event, but workflow executions can still be created manually for a specific realm resource.

At the moment, workflows support event functions for the following realm resources:

#### [](#_workflow_event_functions_)Event functions

   Event Description Parameters

`user-created`

User is added to the realm.

None

`user-authenticated`

User authenticates to the realm.

The client id of the client acting on behalf of the user when authenticating.

`user-federated-identity-added`

User is federated or when the account is linked to an identity provider.

The identity provider alias.

`user-federated-identity-removed`

User’s federated identity is removed or unlinked from an identity provider.

The identity provider alias.

`user-role-granted`

A role is granted to the user.

The role name.

`user-role-revoked`

A role is revoked from the user.

The role name.

`user-group-membership-added`

User joined a group.

The group name or path

`user-group-membership-removed`

User is removed from a group.

The group name or path

`client-created`

Client is added to the realm.

None

`client-authenticated`

Client authenticates to the realm.

None

### [](#_scheduling_workflows_)Scheduling workflows

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/scheduling-workflows.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Fscheduling-workflows.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Fscheduling-workflows.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Workflows can be scheduled to run periodically according to a defined interval. This is done using the `schedule` setting in the workflow definition.

A key difference between event-based triggering and scheduled workflows is that scheduled workflows follows a passive execution model where the workflow engine periodically checks for realm resources that meet the defined condition.

```
name: Track inactive users
schedule:
  after: 5s
  batch-size: 100
```

This method of scheduling is useful for automating tasks that need to be performed regularly, such as cleaning up inactive users or enforcing specific policies on realm resources. It is an alternative to event-based triggering, but it can also be used in combination with it. When used together, the workflow will be triggered either by the defined event or by the schedule.

When a workflow is scheduled, it will be triggered automatically at the defined interval. At each run, the workflow engine will query for the realm resources that matches the workflow’s condition and will create a workflow execution for each of them, up to the defined batch size. If no condition is defined, the workflow will be executed for all realm resources of the type associated with the workflow. In the example above, the workflow is scheduled to run every 5 seconds and will process up to 100 realm resources at each run.

The `schedule` setting supports the following parameters:

- `after`: Defines the interval between each run of the workflow.
- `batch-size`: Defines the maximum number of realm resources to process at each run.

### [](#_workflow_conditions_)Defining conditions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/defining-conditions.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Fdefining-conditions.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Fdefining-conditions.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The optional `if` setting allows you to define the conditions, as expressions, that the target resource must meet in order for the workflow to be triggered. See [Understanding the workflow expression language](#_workflow_expression_language_) for more details.

Conditions provide fine-grained control over whether a workflow execution should be created. They allow you to inspect the context of the event and the state of the resource. For example, you can create conditions to check: * If a user has a specific attribute. * If a specific role is granted to a user or if the user is joining a specific group.

If the condition evaluates to `true`, the workflow execution is created. If it evaluates to `false`, no workflow execution is created, even though the expression from the `on` setting evaluates to `true`.

Just like the `on` setting, the condition is written using an expression that supports a variety of checks on the realm resource associated with the event. For instance, considering a `user_created` event, you can define a condition to trigger the workflow only if the user has a specific attribute:

```
on: user-created
if: has-user-attribute(plan=gold)
```

In this example, the workflow will only be triggered when a new user is created and that user has an attribute `plan` with the value `gold`.

Keycloak provides a set of built-in conditions that you can use in your workflows. The conditions are also based on the realm resource associated with the event.

#### [](#_workflow_user_functions_)User functions

   Condition Description Parameters

`has-user-attribute`

If the user has an attribute set.

The attribute name and optionally the attribute value using a properties format. If multiple values, they should be separated by comma. If the value is omitted, only the presence of the attribute is checked.

`has-role`

If the user is granted with a specific role

The name of the role.

`has-identity-provider-link`

If the user is linked to an identity provider.

The alias of the identity provider.

`is-member-of`

If the user is member of a specific group.

The name or path of the group.

### [](#_workflow_steps_)Defining steps

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/defining-steps.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Fdefining-steps.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Fdefining-steps.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The `steps` setting allows you to define the step chain, which is a sequence of actions to be executed during the lifetime of a workflow execution. Each step represents a specific action that can be performed, such as sending a notification, updating user attributes, or interacting with external systems.

Each step in the step chain is defined using the following structure:

```
steps:
  - uses: <step>
    with:
      <param>: <value>
      ...
    after: <duration>
  ...
```

In its simplest form, a step is defined using only the `uses` field:

```
steps:
  - uses: disable-user
```

However, some steps also support additional settings to customize their behavior:

```
steps:
  - uses: notify-user
    with:
      message: "Welcome to the platform!"
```

Here is a more detailed look at all fields available from defining a step:

`uses`

The name of the step to be executed. This field is mandatory.

`with`

Key/value pairs to be passed to the step. The available parameters depend on the step being used. This field is optional.

`after`

An optional duration to wait before executing the step. The duration is defined using a number followed by a time unit:

- `ms`: milliseconds
- `s`: seconds (default if no unit is specified)
- `m`: minutes
- `h`: hours
- `d`: days

In addition, ISO-8601 duration format is also supported, for example: `P1DT2H` (1 day and 2 hours). If no time unit is specified, it assumes seconds by default.

Keycloak provides a set of built-in steps that you can use in your workflows. The steps are targeted at performing actions on the realm resource associated with the event, so that each realm resource type has its own set of steps.

#### [](#_workflow_user_steps_)User steps

   Step Description Configuration

`add-required-action`

Add a required action to the user

- `action`: The name of the required action

`remove-required-action`

Remove a required action from the user

- `action`: The name of the required action

`grant-role`

Grant one or more roles to the user

- `role`: One or more role names to grant. This can be a single value or a list of values (e.g., `[value1, value2]`)

`revoke-role`

Revoke one or more roles from the user

- `role`: One or more role names to revoke. This can be a single value or a list of values (e.g., `[value1, value2]`)

`join-group`

Add the user to one or more groups

- `group`: One or more group names or paths to join. This can be a single value or a list of values (e.g., `[value1, value2]`)

`leave-group`

Remove the user from one or more groups

- `group`: One or more group names or paths to leave. This can be a single value or a list of values (e.g., `[value1, value2]`)

`set-user-attribute`

Set one or more attributes on the user. Allows providing multiple `<name>`/`<value>` pairs

- `<name>`: The attribute name
- `<value>`: The value of the attribute

`remove-user-attribute`

Remove one or more attributes from the user

- `attribute`: One or more attribute names to remove. This can be a single value or a list of values (e.g., `[value1, value2]`)

`notify-user`

Notify the user by email

- `subject`: The email subject
- `message`: The email message in plain text or HTML format
- `to`: The recipient email address. If not provided, the user’s email address will be used

`unlink-user`

Unlink the user from one or more external Identity Providers

- `idp`: One or more Identity Provider aliases to unlink. This can be:
  
  - Single value
  - List of values (e.g., `[value1, value2]`)
  - `*` to unlink user from all linked Identity Providers

`disable-user`

Disable the user

None

`delete-user`

Delete the user

None

#### [](#_workflow_client_steps_)Client steps

   Step Description Configuration

`delete-client`

Delete the client

None

`disable-client`

Disable the client

None

#### [](#_workflow_immediate_steps_)Understanding immediate steps

A step can run immediately whenever it is reached in the step chain. A immediate step is executed as soon as the previous step is completed, without any delay. This is the default behavior for steps in a workflow.

For example, in the following workflow definition, both steps are immediate steps:

```
steps:
  - uses: notify-user
  - uses: disable-user
```

#### [](#understanding-scheduled-steps)Understanding scheduled steps

A step can be scheduled to run after a specific duration using the `after` field. A scheduled step is executed after waiting for the defined duration once the previous step is completed.

For example, in the following workflow definition, the `disable-user` step is a scheduled step that will run 7 days after notifying the user:

```
steps:
  - uses: notify-user
  - uses: disable-user
    after: 7d
```

### [](#_understanding_workflows_engine_)Understanding the workflows engine

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/understanding-workflows-engine.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-workflows-engine.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-workflows-engine.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The lifecycle and execution of workflows in Keycloak is managed by the Workflows Engine. The engine is responsible for processing events, creating workflow executions, and processing their steps. Once a workflow execution is created, the engine takes over and manages the execution of its steps chain.

The execution of the steps chain is done asynchronously and detached from the request that originated the event that triggered the workflow. This means that workflows are not scheduled nor executed immediately within the context of the request, but rather sent to a task thread pool and scheduled to be executed later by the engine. This allows for better performance and scalability, as the steps can be executed in the background without blocking the main request thread.

For immediate steps, the engine processes them as soon as the workflow execution is submitted to the task thread pool. This means that immediate steps are executed right away, as soon as the workflow execution is picked up by a worker thread from the pool.

For scheduled steps, the engine continuously monitors the state of all active workflow executions in the realm. This is done by a background task that runs periodically to check for any workflow executions that have steps ready to be executed. The schedule steps of a workflow are also executed asynchronously by the background task, using the same task thread pool as immediate steps.

#### [](#configuring-the-scheduled-steps-execution-interval)Configuring the scheduled steps execution interval

By default, the background task tracking due steps runs every 12 hours, but this interval can be configured by setting the following server configuration option:

- `spi-events-listener—​workflow-event-listener—​step-runner-task-interval`: Defines the interval at which the background task runs to check for workflow steps that are due for execution. It follows the same format used in the workflow step `after` field, where you can specify the interval as a number followed by a time unit (`ms`, `s` - default, `m`, `h`, `d`) or using the ISO-8601 duration format.

By default, the first execution of the background task occurs after one full interval has elapsed since the server started. You can control this by setting a start time that acts as an anchor for aligning executions to a predictable schedule:

- `spi-events-listener—​workflow-event-listener—​step-runner-task-start-time`: Defines the time of day used to align the background task executions, in `HH:mm` format (e.g., `02:00`, `14:30`), using the server’s local timezone.

When a start time is set, the task interval is aligned to a grid of execution times anchored at the specified time. For example, with a start time of `02:00` and an interval of `12h`, the task always runs at 02:00 and 14:00 regardless of when the server was started. If the server starts at 10:30, the first execution would occur at 14:00.

You can adjust these options based on your realm’s needs and the expected frequency of workflow executions.

#### [](#configuring-the-task-execution-timeout)Configuring the task execution timeout

Most of the time, steps are executed in a short amount of time. However, there might be cases where a step takes longer to execute due to various reasons, such as network latency depending on an external service delays, or complex processing logic.

By default, the timeout for executing a workflow is set to **5 seconds**. The amount of steps that are processed in a single run depends on the number of immediate steps that are defined in the workflow, as they are executed sequentially in the same run. Once the workflow encounters a scheduled step, it schedules the step and terminates the execution.

If a workflow execution takes longer than the defined timeout, the engine will fail the current step, and it will attempt to run this step in the next execution of the workflow. This is done to prevent long-running executions from blocking the workflow task thread pool and impacting the execution of other workflows.

You can configure the timeout for the workflow executor task by setting the following server configuration option:

- `spi-workflow—​default—​executor-task-timeout`: Defines the timeout for executing a workflow. It follows the same format used in the workflow step `after` field, where you can specify the interval as a number followed by a time unit (`ms`, `s` - default, `m`, `h`, `d`) or using the ISO-8601 duration format.

#### [](#performance-considerations-2)Performance considerations

Workflows can introduce additional processing overhead. This is mainly true when workflows are triggered frequently such as when users are authenticating. Even though workflows are scheduled and executed asynchronously, they still consume system resources. Therefore, it’s important to monitor the performance of your realm and adjust your workflow definitions accordingly to ensure optimal performance.

Consider the following best practices when defining workflows:

- Keep workflows simple and focused on specific tasks and avoid long-running transactions when executing steps in a workflow. This is specially true for steps that are executed immediately upon workflow execution creation.
- Consider using short timeouts for the step execution to avoid blocking the workflow execution for too long and potentially exhaust the workflow task thread pool and impacting the execution of other workflows.
- Adjust the background task interval based on the expected frequency of workflow executions in your realm. If your realm has a high volume of workflow executions, consider reducing the interval to avoid processing multiple workflows at once.

### [](#_handling_failures_)Handling failures

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/handling-failures.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Fhandling-failures.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Fhandling-failures.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The workflows engine keeps track of the execution process by storing the step that should run in a state table. If the step fails to run, either due to an error in the step execution or because of a timeout, the error is logged, an event is fired, and the state table remains unchanged. This effectively means that the step will be retried the next time the workflow execution task runs.

In this initial version there’s no limit to the number of retries, so a workflow execution can get stuck until the administrator intervenes and either fixes the issue that is preventing the step from running successfully or uses the API to cancel the workflow execution or to migrate the resource to a different workflow/step. Thus, it is important that admins monitor the workflow execution logs and check for any errors that may occur repeatedly.

The state table is used even for immediate steps (i.e. steps that are supposed to run immediately after the previous step). This means that if an immediate step fails, the workflow execution will be retried later, and the failed step will be retried as well, behaving as if it were a scheduled step. This is to ensure that the workflow execution process is consistent and that all steps are retried in the same way, regardless of their configuration. This also ensures the workflow will be resumed in case of server restarts or crashes.

Future versions of the workflows engine will include more features to handle failures, such as the ability to configure a maximum number of retries for each step, as well as the ability to define custom error handling logic for specific steps, like skip the step or cancel the workflow execution.

### [](#_troubleshooting_workflows_)Troubleshooting workflows

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/troubleshooting-workflows.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Ftroubleshooting-workflows.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Ftroubleshooting-workflows.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When workflows are not behaving as expected, server logs are the primary tool for understanding what is happening. The workflow engine logs detailed information about workflow lifecycle events, step execution, and errors.

#### [](#enabling-workflow-debug-logging)Enabling workflow debug logging

By default, the workflow engine logs warnings and errors. To see the full execution flow, enable `DEBUG` level logging for the `org.keycloak.models.workflow` category by setting the following server configuration option:

```
--log-level=org.keycloak.models.workflow:debug
```

This enables debug logging for the workflow engine, all step providers, and the workflow state provider.

#### [](#what-to-look-for-in-the-logs)What to look for in the logs

##### [](#workflow-activation-and-scheduling)Workflow activation and scheduling

When a workflow is triggered by an event, the engine logs the activation and scheduling of the first step:

- `Workflow '<name>' activated for resource <id> (execution id: <id>)` — the workflow execution has been created and is being processed.
- `Scheduled step <step> to run in <duration> for resource <id> (execution id: <id>)` — a scheduled step has been queued to run after the specified duration.

##### [](#step-execution)Step execution

As steps are executed, the engine logs their progress:

- `Running step <step> on resource <id> (execution id: <id>)` — a step is about to be executed.
- `Step <step> completed successfully (execution id: <id>)` — the step finished without errors.
- `Workflow '<name>' completed for resource <id> (execution id: <id>)` — all steps in the workflow have been executed.

##### [](#step-failures-and-timeouts)Step failures and timeouts

When a step fails or the executor times out, the engine logs the error and fires a step failed event. The step will be retried in the next execution cycle (see [Handling failures](#_handling_failures_) for more details).

- `Step <step> failed (execution id: <id>) - error message: <message>` — the step threw an exception during execution.
- `Workflow executor timed out during execution of step <step>` — runtime exception indicating that the step exceeded the configured executor task timeout. See [Understanding the workflows engine](#_understanding_workflows_engine_) for how to configure the timeout.
- `Workflow executor was cancelled during execution of step <step>` — runtime exception indicating that the executor service was cancelled for a reason other than timeout.

##### [](#skipped-and-cancelled-workflows)Skipped and cancelled workflows

- `Skipping workflow <name> as it is disabled` — the workflow is disabled and will not process events.
- `Workflow '<name>' cancelled for resource <id> (execution id: <id>)` — the workflow execution was cancelled.
- `Resource <id> is no longer eligible for workflow <name>. Cancelling execution of the workflow.` — the resource no longer meets the workflow conditions upon resuming a scheduled step.

##### [](#workflow-errors)Workflow errors

- `Error processing event <event> for workflow <name>: <message>` — an error occurred while processing the event that triggers the workflow.
- `Error resuming workflow <name> for resource <id>: <message>` — an error occurred while resuming a scheduled workflow execution.
- `Could not find step <step> in workflow <name> for resource <id>. Cancelling execution of the workflow.` — the step scheduled to run could not be found in the workflow.

##### [](#step-level-warnings)Step-level warnings

Individual step providers log warnings when configuration is missing or invalid. Examples include:

- `Missing required configuration option '<key>' in <step>` — a required parameter was not provided in the step’s `with` block.
- `Invalid required action <action> configured in <step>` — the required action name is not recognized.
- `User <id> has no email address, skipping notification` — the `notify-user` step could not send an email.
- `Failed to send notification email to <recipient>` — an email delivery failure occurred.

#### [](#useful-log-categories)Useful log categories

If you need more targeted logging, you can enable debug logging for specific components:

  Category Description

`org.keycloak.models.workflow`

All workflow-related logging (engine, steps, state)

`org.keycloak.models.workflow.DefaultWorkflowProvider`

Workflow engine lifecycle (activation, conditions, event processing)

`org.keycloak.models.workflow.RunWorkflowTask`

Step execution, completion, failures, and timeouts

`org.keycloak.models.workflow.ScheduleWorkflowTask`

Step scheduling and delayed execution

### [](#_understanding_common_use_cases_)Understanding common use cases

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/workflows/understanding-common-use-cases.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-common-use-cases.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fworkflows%2Funderstanding-common-use-cases.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Workflows can be used to automate a wide range of administrative tasks within a realm. Here are some common use cases where workflows can be particularly beneficial:

#### [](#user-onboarding)User Onboarding

When a new user is created, a workflow can automatically send a welcome email, and assign initial roles or groups to the user. This ensures that new users have the necessary access and information from the moment they join.

```
name: Onboarding gold members
on: user-created
if: has-user-attribute(membership=gold)
steps:
  - uses: notify-user
    with:
      message: "Welcome to the Gold Membership program!"
  - uses: grant-role
    with:
      role: gold-member
```

#### [](#user-offboarding)User Offboarding

When a user is removed from a specific group, a workflow can automatically revoke certain roles or permissions associated with that group.

```
name: Offboarding sales members
on: user-group-membership-removed(/Sales)
steps:
  - uses: revoke-role
    with:
      role:
        - sales-rep
        - manager
        - sales-intern
```

#### [](#tracking-user-activity-and-taking-actions-on-inactivity)Tracking user activity and taking actions on inactivity

When a user has been inactive for a certain period, a workflow can send reminder emails and deactivate the account.

```
name: Track inactive users
on: user-authenticated
schedule:
  after: 5s
  batch-size: 2
concurrency:
  restart-in-progress: true
steps:
  - uses: notify-user
    after: 180d
    with:
      message: It has been a while since your last login. We miss you!
  - uses: notify-user
    after: 60d
    with:
      message: Your account will be disabled in ${workflow.daysUntilNextStep} days!
  - uses: disable-user
    after: 7d
  - uses: notify-user
    with:
      message: Your account was disabled. Sorry to see you go.
```

## [](#assembly-managing-clients_server_administration_guide)Managing OpenID Connect and SAML Clients

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/assembly-managing-clients.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fassembly-managing-clients.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fassembly-managing-clients.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Clients are entities that can request authentication of a user. Clients come in two forms. The first type of client is an application that wants to participate in single sign-on. These clients just want Keycloak to provide security for them. The other type of client is one that is requesting an access token so that it can invoke other services on behalf of the authenticated user. This section discusses various aspects around configuring clients and various ways to do it.

### [](#_oidc_clients)Managing OpenID Connect clients

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/assembly-client-oidc.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fassembly-client-oidc.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fassembly-client-oidc.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

[OpenID Connect](#con-oidc_server_administration_guide) is the recommended protocol to secure applications. It was designed from the ground up to be web friendly and it works best with HTML5/JavaScript applications.

#### [](#proc-creating-oidc-client_server_administration_guide)Creating an OpenID Connect client

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/proc-creating-oidc-client.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fproc-creating-oidc-client.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fproc-creating-oidc-client.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To protect an application that uses the OpenID connect protocol, you create a client.

Procedure

1. Click **Clients** in the menu.
2. Click **Create client**.
   
   Create client
   
   ![Create Client](./images/add-client-oidc.png)
3. Leave **Client type** set to **OpenID Connect**.
4. Enter a **Client ID.**
   
   This ID is an alphanumeric string that is used in OIDC requests and in the Keycloak database to identify the client.
5. Supply a **Name** for the client.
   
   If you plan to localize this name, set up a replacement string value. For example, a string value such as ${myapp}. See the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more information.
6. Click **Save**.

This action creates the client and bring you to the **Settings** tab, where you can perform [Basic configuration](#con-basic-settings_server_administration_guide).

#### [](#con-basic-settings_server_administration_guide)Basic configuration

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/con-basic-settings.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-basic-settings.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-basic-settings.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The **Settings** tab includes many options to configure this client.

Settings tab

![Settings tab](./images/client-settings-oidc.png)

##### [](#general-settings)General Settings

**Client ID**

The alphanumeric ID string that is used in OIDC requests and in the Keycloak database to identify the client.

**Name**

The name for the client in Keycloak UI screen. To localize the name, set up a replacement string value. For example, a string value such as ${myapp}. See the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more information.

**Description**

The description of the client. This setting can also be localized.

**Always Display in Console**

Always list this client in the Account Console even if this user does not have an active session.

##### [](#access-settings)Access Settings

**Root URL**

If Keycloak uses any configured relative URLs, this value is prepended to them.

**Home URL**

Provides the default URL for when the auth server needs to redirect or link back to the client.

**Valid Redirect URIs**

Required field. Enter a URL pattern and click **+** to add and **-** to remove existing URLs and click **Save**. Exact (case sensitive) string matching is used to compare valid redirect URIs.

You can use wildcards at the end of the URL pattern. For example `http://host.com/path/*`. To avoid security issues, if the passed redirect URI contains the **userinfo** part or its **path** manages access to parent directory (`/../`) no wildcard comparison is performed but the standard and secure exact string matching.

The full wildcard `*` valid redirect URI can also be configured to allow any **http** or **https** redirect URI. Please do not use it in production environments.

Exclusive redirect URI patterns are typically more secure. See [Unspecific Redirect URIs](#unspecific-redirect-uris_server_administration_guide) for more information.

Web Origins

Enter a URL pattern and click + to add and - to remove existing URLs. Click Save.

This option handles [Cross-Origin Resource Sharing (CORS)](https://fetch.spec.whatwg.org/). If browser JavaScript attempts an AJAX HTTP request to a server whose domain is different from the one that the JavaScript code came from, the request must use CORS. The server must handle CORS requests, otherwise the browser will not display or allow the request to be processed. This protocol protects against XSS, CSRF, and other JavaScript-based attacks.

Domain URLs listed here are embedded within the access token sent to the client application. The client application uses this information to decide whether to allow a CORS request to be invoked on it. Only Keycloak client adapters support this feature. See [Securing applications Guides](https://www.keycloak.org/guides#securing-apps) for more information.

Admin URL

Callback endpoint for a client. The server uses this URL to make callbacks like pushing revocation policies, performing backchannel logout, and other administrative operations. For Keycloak servlet adapters, this URL can be the root URL of the servlet application. The callback messages sent to this URL are sent in the Keycloak specific format, which is not OIDC standard. This format is supported only for clients secured by the legacy Keycloak Java OIDC adapters or by the [Elytron Wildfly OIDC adapter](https://docs.wildfly.org/37/WildFly_Elytron_Security.html#Keycloak_Integration).

##### [](#capability-config)Capability Config

**Client authentication**

Specifies the type of OIDC client.

- *ON*
  
  For server-side clients that perform browser logins and require client secrets when making an Access Token Request. This setting should be used for server-side applications. Clients with the client authentication enabled are referred as confidential clients. For more details, see [Client credentials](#_client-credentials).
- *OFF*
  
  For client-side clients that perform browser logins. As it is not possible to ensure that secrets can be kept safe with client-side clients, it is important to restrict access by configuring correct redirect URIs. Clients with the client authentication disabled are referred as public clients.

**Authorization**

Enables or disables fine-grained authorization support for this client.

**Standard Flow**

If enabled, this client can use the OIDC [Authorization Code Flow](#_oidc-auth-flows-authorization).

**Direct Access Grants**

If enabled, this client can use the OIDC [Direct Access Grants](#_oidc-auth-flows-direct).

**Implicit Flow**

If enabled, this client can use the OIDC [Implicit Flow](#_oidc-auth-flows-implicit).

**Service account roles**

If enabled, this client can authenticate to Keycloak and retrieve access token dedicated to this client. In terms of OAuth2 specification, this enables support of `Client Credentials Grant` for this client.

**Standard Token Exchange**

If enabled, this client can use the [Standard token exchange](https://www.keycloak.org/securing-apps/token-exchange#_standard-token-exchange).

**Auth 2.0 Device Authorization Grant**

If enabled, this client can use the OIDC [Device Authorization Grant](#con-oidc-auth-flows_server_administration_guide).

**OIDC CIBA Grant**

If enabled, this client can use the OIDC [Client Initiated Backchannel Authentication Grant](#con-oidc-auth-flows_server_administration_guide).

**PKCE method**

If an attacker steals an authorization code of a legitimate client, Proof Key for Code Exchange (PKCE) prevents the attacker from receiving the tokens that apply to the code. With this option, you can specify which PKCE challenge method is required for this client.

An administrator can select one of these options:

**(blank)**

Keycloak does not apply PKCE unless the client sends the appropriate PKCE parameters to Keycloak authorization endpoint. So PKCE is still possible to use, but it is not required.

**S256**

Keycloak applies to the client PKCE whose code challenge method is S256.

**plain**

Keycloak applies to the client PKCE whose code challenge method is plain.

See [RFC 7636 Proof Key for Code Exchange by OAuth Public Clients](https://datatracker.ietf.org/doc/html/rfc7636) for more details.

**Require DPoP bound tokens**

DPoP binds an access token and a refresh token together with the public part of a client’s key pair. For the details, see [DPoP](#_dpop-bound-tokens).

##### [](#login-settings)Login settings

**Login theme**

A theme to use for login, OTP, grant registration, and forgotten password pages.

**Consent required**

If enabled, users have to consent to client access.

For client-side clients that perform browser logins. As it is not possible to ensure that secrets can be kept safe with client-side clients, it is important to restrict access by configuring correct redirect URIs.

**Display client on screen**

This switch applies if **Consent Required** is **Off**.

- *Off*
  
  The consent screen will contain only the consents corresponding to configured client scopes.
- *On*
  
  There will be also one item on the consent screen about this client itself.

**Client consent screen text**

Applies if **Consent required** and **Display client on screen** are enabled. Contains the text that will be on the consent screen about permissions for this client.

##### [](#logout-settings)Logout settings

**Front channel logout**

If **Front Channel Logout** is enabled, the application should be able to log out users through the front channel as per [OpenID Connect Front-Channel Logout](https://openid.net/specs/openid-connect-frontchannel-1_0.html) specification. If enabled, you should also provide the `Front-Channel Logout URL`.

**Front-channel logout URL**

URL that will be used by Keycloak to send logout requests to clients through the front-channel. If not provided, it defaults to the Home URL. This option is applicable just if `Front channel logout` option is ON.

**Front-channel logout session required**

Specifies whether a sid (session ID) and iss (issuer) parameters are included in the Logout request when the Front-channel Logout URL is used.

**Backchannel logout URL**

URL that will cause the client to log itself out when a logout request is sent to this realm (via end\_session\_endpoint). The logout is done by sending logout token as specified in the OIDC Backchannel logout specification. If omitted, the logout request might be sent to the specified `Admin URL` (if configured) in the format specific to Keycloak adapters. If even `Admin URL` is not configured, no logout request will be sent to the client. This option is applicable just if `Front channel logout` option is OFF.

**Backchannel logout session required**

Specifies whether a session ID Claim is included in the Logout Token when the **Backchannel Logout URL** is used.

**Backchannel logout revoke offline sessions**

Specifies whether a revoke\_offline\_access event is included in the Logout Token when the Backchannel Logout URL is used. Keycloak will revoke offline sessions when receiving a Logout Token with this event.

**Logout confirmation**

When enabled, Keycloak displays a confirmation page to the user after a successful logout that reads “You are logged out”. This setting primarily affects browser-based logouts, including [OIDC Logout](#_oidc-logout) initiated by the client (RP-Initiated Logout). If a `post_logout_redirect_uri` is provided and validated for this client, the confirmation page includes a link (or button) to continue to that URL instead of redirecting automatically.

#### [](#con-advanced-settings_server_administration_guide)Advanced configuration

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/con-advanced-settings.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-advanced-settings.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-advanced-settings.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

After completing the fields on the **Settings** tab, you can use the other tabs to perform advanced configuration. For example, you can use the **Roles** or **Client scopes** tabs to configure client roles defined for the client or manage client scopes for the client. Also, see the remaining sections in this chapter for other capabilities.

##### [](#advanced-tab)Advanced tab

When you click the **Advanced** tab, additional fields are displayed. For details on a specific field, click the question mark icon for that field. However, certain fields are described in detail in this section.

##### [](#fine-grain-openid-connect-configuration)Fine grain OpenID Connect configuration

**Logo URL**

URL that references a logo for the Client application.

**Policy URL**

URL that the Relying Party Client provides to the End-User to read about how the profile data will be used.

**Terms of Service URL**

URL that the Relying Party Client provides to the End-User to read about the Relying Party’s terms of service.

**Signed and Encrypted ID Token Support**

Keycloak can encrypt ID tokens according to the [Json Web Encryption (JWE)](https://datatracker.ietf.org/doc/html/rfc7516) specification. The administrator determines if ID tokens are encrypted for each client.

The key used for encrypting the ID token is the Content Encryption Key (CEK). Keycloak and a client must negotiate which CEK is used and how it is delivered. The method used to determine the CEK is the Key Management Mode. The Key Management Mode that Keycloak supports is Key Encryption.

In Key Encryption:

1. The client generates an asymmetric cryptographic key pair.
2. The public key is used to encrypt the CEK.
3. Keycloak generates a CEK per ID token
4. Keycloak encrypts the ID token using this generated CEK
5. Keycloak encrypts the CEK using the client’s public key.
6. The client decrypts this encrypted CEK using their private key
7. The client decrypts the ID token using the decrypted CEK.

No party, other than the client, can decrypt the ID token.

The client must pass its public key for encrypting CEK to Keycloak. Keycloak supports downloading public keys from a URL provided by the client. The client must provide public keys according to the [Json Web Keys (JWK)](https://datatracker.ietf.org/doc/html/rfc7517) specification.

The procedure is:

1. Open the client’s **Keys** tab.
2. Toggle **JWKS URL** to ON.
3. Input the client’s public key URL in the **JWKS URL** textbox.

Key Encryption’s algorithms are defined in the [Json Web Algorithm (JWA)](https://datatracker.ietf.org/doc/html/rfc7518#section-4.1) specification. Keycloak supports:

- RSAES-PKCS1-v1\_5(RSA1\_5)
- RSAES OAEP using default parameters (RSA-OAEP)
- RSAES OAEP 256 using SHA-256 and MFG1 (RSA-OAEP-256)

The procedure to select the algorithm is:

1. Open the client’s **Advanced** tab.
2. Open **Fine Grain OpenID Connect Configuration**.
3. Select the algorithm from **ID Token Encryption Content Encryption Algorithm** pulldown menu.

##### [](#openid-connect-compatibility-modes)OpenID Connect Compatibility Modes

This section exists for backward compatibility. Click the question mark icons for details on each field.

**OAuth 2.0 Mutual TLS Certificate Bound Access Tokens Enabled**

Mutual TLS binds an access token and a refresh token together with a client certificate, which is exchanged during a TLS handshake. This binding prevents an attacker from using stolen tokens.

This type of token is a holder-of-key token. Unlike bearer tokens, the recipient of a holder-of-key token can verify if the sender of the token is legitimate.

If this setting is on, the workflow is:

1. A token request is sent to the token endpoint in an authorization code flow or hybrid flow.
2. Keycloak requests a client certificate.
3. Keycloak receives the client certificate.
4. Keycloak successfully verifies the client certificate.

If verification fails, Keycloak rejects the token.

In the following cases, Keycloak will verify the client sending the access token or the refresh token:

- A token refresh request is sent to the token endpoint with a holder-of-key refresh token.
- A UserInfo request is sent to UserInfo endpoint with a holder-of-key access token.
- A logout request is sent to non-OIDC compliant Keycloak proprietary Logout endpoint with a holder-of-key refresh token.

See [Mutual TLS Client Certificate Bound Access Tokens](https://datatracker.ietf.org/doc/rfc8705/) in the OAuth 2.0 Mutual TLS Client Authentication and Certificate Bound Access Tokens for more details.

Sender-constrained tokens like X.509 certificate-bound tokens (mTLS tokens with the `x5t#S256` confirmation claim), cannot be used as the `subject_token` parameter in [Standard Token Exchange](https://www.keycloak.org/securing-apps/token-exchange#_standard-token-exchange).

Keycloak client adapters do not support holder-of-key token verification. Keycloak adapters treat access and refresh tokens as bearer tokens.

**Advanced Settings for OIDC**

The Advanced Settings for OpenID Connect allows you to configure overrides at the client level for [session and token timeouts](#_timeouts).

![Advanced Settings](./images/client-advanced-settings-oidc.png)

  Configuration Description

Access Token Lifespan

The value overrides the realm option with same name.

Client Session Idle

The value overrides the realm option with same name. The value should be shorter than the global **SSO Session Idle**.

Client Session Max

The value overrides the realm option with same name. The value should be shorter than the global **SSO Session Max**.

Client Offline Session Idle

This setting allows you to configure a shorter offline session idle timeout for the client. The timeout is amount of time the session remains idle before Keycloak revokes its offline token. If not set, realm [Offline Session Idle](#_offline-session-idle) is used.

Client Offline Session Max

This setting allows you to configure a shorter offline session max lifespan for the client. The lifespan is the maximum time before Keycloak revokes the corresponding offline token. This option needs [Offline Session Max Limited](#_offline-session-max-limited) enabled globally in the realm, and defaults to [Offline Session Max](#_offline-session-max).

**ACR to Level of Authentication (LoA) Mapping**

In the advanced settings of a client, you can define which `Authentication Context Class Reference (ACR)` value is mapped to which `Level of Authentication (LoA)`. This mapping can be specified also at the realm as mentioned in the [ACR to LoA Mapping](#_mapping-acr-to-loa-realm). A best practice is to configure this mapping at the realm level, which allows to share the same settings across multiple clients.

The `Default ACR Values` can be used to specify the default values when the login request is sent from this client to Keycloak without `acr_values` parameter and without a `claims` parameter that has an `acr` claim attached. See [official OIDC dynamic client registration specification](https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata).

Note that default ACR values are used as the default level, however it cannot be reliably used to enforce login with the particular level. For example, assume that you configure the `Default ACR Values` to level 2. Then by default, users will be required to authenticate with level 2. However, when the user explicitly attaches the parameter into login request such as `acr_values=1`, then the level 1 will be used. As a result, if the client really requires level 2, the client is encouraged to check the presence of the `acr` claim inside ID Token and double-check that it contains the requested level 2. To actually enforce the usage of a certain ACR on the Keycloak side, use the `Minimum ACR Value` setting. This allows administrators to enforce ACRs even on applications that are not able to validate the requested `acr` claim inside the token.

![ACR to LoA mapping](./images/client-oidc-map-acr-to-loa.png)

For further details see [Step-up Authentication](#_step-up-flow) and [the official OIDC specification](https://openid.net/specs/openid-connect-core-1_0.html#acrSemantics).

#### [](#_client-credentials)Confidential client credentials

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/con-confidential-client-credentials.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-confidential-client-credentials.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-confidential-client-credentials.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

If the [Client authentication](#_access-type) of the client is set to **ON**, the credentials of the client must be configured under the **Credentials** tab.

Credentials tab

![Credentials Tab](./images/client-credentials.png)

The **Client Authenticator** drop-down list specifies the type of credential to use for your client.

**Client ID and Secret**

This choice is the default setting. A random secret is generated automatically, but you can override it with a specific value if needed. The client secret can also reference a value stored in an [external vault](#_vault-administration). Click **Regenerate** to recreate random secret if necessary.

As specified in the OAuth2 and OpenID Connect specifications, it is possible to authenticate client either with the client secret in the `Authorization: Basic` header (together with `client_id`) or by sending `client_id` and `client_secret` as parameters in the HTTP POST method body. You can configure the **Allowed authentication method** to restrict only to one specified method of those if you want to authenticate client for example just with the `Authorization: Basic` method, but never with the request parameters sent in the HTTP POST method body.

**Signed JWT issued by the client** ![Signed JWT](./images/client-credentials-jwt.png)

**Signed JWT** allows a client to authenticate with self-signed client assertions. This enables the client to authenticate without a shared secret.

In this authenticator you can enforce the **Signature algorithm** used by the client (any algorithm is valid by default) and the **Max expiration** allowed for the JWT token (tokens received after this period will not be accepted because they are too old, note that tokens should be issued right before the authentication, 60 seconds by default).

When choosing this credential type you will have to also generate a private key and certificate for the client in the tab `Keys`. The private key will be used to sign the JWT, while the certificate is used by the server to verify the signature.

Keys tab

![Keys tab](./images/client-oidc-keys.png)

Click on the `Generate new keys` button to start this process.

Generate keys

![generate client keys](./images/generate-client-keys.png)

1. Select the archive format you want to use.
2. Enter a **key password**.
3. Enter a **store password**.
4. Click **Generate**.

When you generate the keys, Keycloak will store the certificate and you download the private key and certificate for your client.

You can also generate keys using an external tool and then import the client’s certificate by clicking **Import Certificate**.

Import certificate

![Import Certificate](./images/import-client-cert.png)

1. Select the archive format of the certificate.
2. Enter the store password.
3. Select the certificate file by clicking **Import File**.
4. Click **Import**.

Importing a certificate is unnecessary if you click **Use JWKS URL**. In this case, you can provide the URL where the public key is published in [JWK](https://datatracker.ietf.org/doc/html/rfc7517) format. With this option, if the key is ever changed, Keycloak reimports the key.

If you are using a client secured by Keycloak adapter, you can configure the JWKS URL in this format, assuming that [https://myhost.com/myapp](https://myhost.com/myapp) is the root URL of your client application:

```
https://myhost.com/myapp/k_jwks
```

See [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more details.

**Signed JWT issued by an Identity Provider**

**Signed JWT** allows a client to authenticate with client assertions issued by an identity provider. Example use-cases include:

- Client assertion issued by an OpenID Connect provider
- SPIFFE JWT SVIDs
- Kubernetes service accounts

Before using this authentication mechanism, an identity provider capable of verifying client assertions should be configured.

The identity providers which support client assertions are:

- [OpenID Connect](#_identity_broker_oidc) (support for client assertions must be enabled)
- [SPIFFE](#_identity_broker_spiffe)

<!--THE END-->

![client federated jwt](./images/client-federated-jwt.png)

- Identity provider - the alias of the identity provider to use
- Federated subject - the external client id for the client (value of the `sub` claim of the client assertion)

**Signed JWT with Client Secret**

If you select this option, you can use a JWT signed by client secret instead of the private key. The client secret can also reference a value stored in an [external vault](#_vault-administration).

The client secret will be used to sign the JWT by the client.

Like in the **Signed JWT** authenticator you can configure the **Signature algorithm** and the **Max expiration** for the JWT token.

**X509 Certificate**

Keycloak will validate if the client uses proper X509 certificate during the TLS Handshake.

X509 certificate

![x509 client auth](./images/x509-client-auth.png)

The validator also checks the Subject DN field of the certificate with a configured regexp validation expression. For some use cases, it is sufficient to accept all certificates. In that case, you can use `(.*?)(?:$)` expression.

Two ways exist for Keycloak to obtain the Client ID from the request:

- The `client_id` parameter in the query (described in Section 2.2 of the [OAuth 2.0 Specification](https://datatracker.ietf.org/doc/html/rfc6749)).
- Supply `client_id` as a form parameter.

#### [](#_dpop-bound-tokens)DPoP

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/con-dpop.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-dpop.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-dpop.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak supports **DPoP** (Demonstrating Proof-of-Possession) to bind access and refresh tokens to a cryptographic key pair, ensuring they can only be used by the legitimate client.

For detailed instructions on how to configure, enforce, and use DPoP within Keycloak, see the dedicated [**Securing applications with DPoP**](https://www.keycloak.org/securing-apps/dpop) guide.

#### [](#_secret_rotation)Client Secret Rotation

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/con-secret-rotation.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-secret-rotation.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-secret-rotation.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Please note that Client Secret Rotation support is in development. Use this feature experimentally.

For a client with [Confidential](#_client-credentials) [Client authentication](#_access-type) Keycloak supports the functionality of rotating client secrets through [Client Policies](#_client_policies).

The client secrets rotation policy provides greater security in order to alleviate problems such as secret leakage. Once enabled, Keycloak supports up to two concurrently active secrets for each client. The policy manages rotations according to the following settings:

- **Secret expiration:** \[seconds] - When the secret is rotated, this is the expiration of time of the new secret. The amount, *in seconds*, added to the secret creation date. Calculated at policy execution time.
- **Rotated secret expiration:** \[seconds] - When the secret is rotated, this value is the remaining expiration time for the old secret. This value should be always smaller than Secret expiration. When the value is 0, the old secret will be immediately removed during client rotation. The amount, *in seconds*, added to the secret rotation date. Calculated at policy execution time.
- **Remaining expiration time for rotation during update:** \[seconds] - Time period when an update to a dynamic client should perform client secret rotation. Calculated at policy execution time.

When a client secret rotation occurs, a new main secret is generated and the old client main secret becomes the secondary secret with a new expiration date.

##### [](#rules-for-client-secret-rotation)Rules for client secret rotation

Rotations do not occur automatically or through a background process. In order to perform the rotation, an update action is required on the client, either through the Keycloak Admin Console through the function of **Regenerate Secret**, in the client’s credentials tab or Admin REST API. When invoking a client update action, secret rotation occurs according to the rules:

- When the value of **Secret expiration** is less than the current date.
- During dynamic client registration client-update request, the client secret will be automatically rotated if the value of **Remaining expiration time for rotation during update** match the period between the current date and the **Secret expiration**.

Additionally it is possible through Admin REST API to force a client secret rotation at any time.

During the creation of new clients, if the client secret rotation policy is active, the behavior will be applied automatically.

To apply the secret rotation behavior to an existing client, update that client after you define the policy so that the behavior is applied.

#### [](#_proc-secret-rotation)Creating an OIDC Client Secret Rotation Policy

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/proc-secret-rotation.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fproc-secret-rotation.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fproc-secret-rotation.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

The following is an example of defining a secret rotation policy:

Procedure

01. Click **Realm Settings** in the menu.
02. Click **Client Policies** tab.
03. On the **Profiles** page, click **Create client profile**.
    
    Create a profile
    
    ![Create Client Profile](./images/create-oidc-client-profile.png)
04. Enter any name for **Name**.
05. Enter a description that helps you identify the purpose of the profile for **Description**.
06. Click **Save**.
    
    This action creates the profile and enables you to configure executors.
07. Click **Add executor** to configure an executor for this profile.
    
    Create a profile executor
    
    ![Client Profile Executor](./images/create-oidc-client-secret-rotation-executor.png)
08. Select *secret-rotation* for **Executor Type**.
09. Enter the maximum duration time of each secret, in seconds, for **Secret Expiration**.
10. Enter the maximum duration time of each rotated secret, in seconds, for **Rotated Secret Expiration**.
    
    Remember that the **Rotated Secret Expiration** value must always be less than **Secret Expiration**.
11. Enter the amount of time, in seconds, after which any update action will update the client for **Remain Expiration Time**.
12. Click **Add**.
    
    In the example above:
    
    - Each secret is valid for one week.
    - The rotated secret expires after two days.
    - The window for updating dynamic clients starts one day before the secret expires.
13. Return to the **Client Policies** tab.
14. Click **Policies**.
15. Click **Create client policy**.
    
    Create the Client Secret Rotation Policy
    
    ![Client Rotation Policy](./images/create-oidc-client-secret-rotation-policy.png)
16. Enter any name for **Name**.
17. Enter a description that helps you identify the purpose of the policy for **Description**.
18. Click **Save**.
    
    This action creates the policy and enables you to associate policies with profiles. It also allows you to configure the conditions for policy execution.
19. Under Conditions, click **Add condition**.
    
    Create the Client Secret Rotation Policy Condition
    
    ![Client Rotation Policy Condition](./images/create-oidc-client-secret-rotation-condition.png)
20. To apply the behavior to all confidential clients select *client-access-type* in the **Condition Type** field
    
    To apply to a specific group of clients, another approach would be to select the *client-roles* type in the **Condition Type** field. In this way, you could create specific roles and assign a custom rotation configuration to each role.
21. Add *confidential* to the field **Client Access Type**.
22. Click **Add**.
23. Back in the policy setting, under *Client Profiles*, click **Add client profile** and then select **Weekly Client Secret Rotation Profile** from the list and then click **Add**.
    
    Client Secret Rotation Policy
    
    ![Client Rotation Policy](./images/oidc-client-secret-rotation-policy.png)

To apply the secret rotation behavior to an existing client, follow the following steps:

Using the Admin Console

1. Click **Clients** in the menu.
2. Click a client.
3. Click the **Credentials** tab.
4. Click **Re-generate** of the client secret.

* * *

Using client REST services it can be executed in two ways:

- Through an update operation on a client
- Through the regenerate client secret endpoint

#### [](#_service_accounts)Using a service account

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/proc-using-a-service-account.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fproc-using-a-service-account.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fproc-using-a-service-account.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Each OIDC client has a built-in *service account*. Use this *service account* to obtain an access token.

Procedure

01. Click **Clients** in the menu.
02. Select your client.
03. Click the **Settings** tab.
04. Toggle [Client authentication](#_access-type) to **On**.
05. Select **Service account roles** checkbox to make sure it is enabled.
06. Click **Save**.
07. Configure your [client credentials](#_client-credentials).
08. Click the **Client Scopes** tab, select the dedicated client scope (usually first client scope in the list, more details [in this section](#_client_scopes_dedicated)) and select **Scope** tab of the client scope.
09. Verify that you have roles or toggle **Full Scope Allowed** to **ON**. Note that this switch is useful only for the development purposes and in the production, it is recommended to disable this switch and properly configure role scopes. The details about this switch are described in [this section](#_role_scope_mappings) and in [this section](#_oidc_token_role_mappings).
10. Click the **Service Account Roles** tab of your client
11. Configure the roles available to this service account for your client.

Roles from access tokens are the intersection of:

- Role scope mappings of a client combined with the role scope mappings inherited from linked client scopes.
- Service account roles.

The REST URL to invoke is `/realms/{realm-name}/protocol/openid-connect/token`. This URL must be invoked as a POST request and requires that you post the client credentials with the request.

By default, client credentials are represented by the clientId and clientSecret of the client in the **Authorization: Basic** header but you can also authenticate the client with a signed JWT assertion or any other custom mechanism for client authentication.

You also need to set the **grant\_type** parameter to "client\_credentials" as per the OAuth2 specification.

For example, the POST invocation to retrieve a service account can look like this:

```
    POST /realms/demo/protocol/openid-connect/token
    Authorization: Basic cHJvZHVjdC1zYS1jbGllbnQ6cGFzc3dvcmQ=
    Content-Type: application/x-www-form-urlencoded

    grant_type=client_credentials
```

Note that the value of `cHJvZHVjdC1zYS1jbGllbnQ6cGFzc3dvcmQ=` used in the `Authorization` header is Base64 encoded value of clientId and clientSecret in the format prescribed by the `Authorization: Basic` header. In this example, the client ID is `product-sa-client` and the client secret was `password` and hence the value was obtained for example by this command in the Unix platform:

```
echo 'product-sa-client:password' | base64
```

Instead of using the header `Authorization: Basic`, it is also possible to send the credentials as parameters `client_id` and `client_secret` of the POST request. For other client credentials methods, the format of the parameters would be different as described above.

The response would be similar to this [Access Token Response](https://datatracker.ietf.org/doc/html/rfc6749#section-4.4.3) from the OAuth 2.0 specification.

```
HTTP/1.1 200 OK
Content-Type: application/json;charset=UTF-8
Cache-Control: no-store
Pragma: no-cache

{
    "access_token":"eyJhbGciOiJSUzI1NiIs...",
    "token_type":"Bearer",
    "expires_in":60,
    "scope": "email profile"
}
```

Only the access token is returned by default. No refresh token is returned and no user session is created on the Keycloak side upon successful authentication by default. Due to the lack of a refresh token, re-authentication is required when the access token expires. However, this situation does not mean any additional overhead for the Keycloak server because sessions are not created by default.

In this situation, logout is unnecessary. However, issued access tokens can be revoked by sending requests to the OAuth2 Revocation Endpoint as described in the [OpenID Connect Endpoints](#con-oidc_server_administration_guide) section.

Additional resources

For more details, see [Client Credentials Grant](#_client_credentials_grant).

#### [](#_oidc_token_role_mappings)Role mappings in the token

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/con-token-role-mappings.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-token-role-mappings.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-token-role-mappings.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

When a user authenticates, there are some roles that are added to the access token. By default, the [Realm roles](#proc-creating-realm-roles_server_administration_guide) are added to the access token into the `realm_access` claim. The [Client roles](#con-client-roles_server_administration_guide) are added by default to the `resource_access` claim.

The roles added to the token are an intersection of:

- Roles, that are [assigned to the user](#_oidc_token_role_mappings_user_roles).
- [Role scope mappings](#_role_scope_mappings) of the roles that the client is permitted to access

##### [](#_oidc_token_role_mappings_user_roles)Roles assigned to the user

Roles assigned to the user can be defined in the Role mappings as described in [this section](#proc-assigning-role-mappings_server_administration_guide). Few details:

- In case that a user is a member of some [groups](#proc-managing-groups_server_administration_guide), then all the roles of these groups are also applied.
- In case that a role is a [composite role](#_composite-roles), the child roles of the composite role are also applied. In the token, the list of the roles is expanded and would contain all the roles.
- In case that the authenticated user is not a normal user, but a [Service account](#_service_accounts), which represents a client, then the service account roles are used. The service account roles are defined in the tab **Service account roles** of the particular client.

##### [](#role-protocol-mappers)Role protocol mappers

Similarly to other claims, the roles are added to the access token issued for the client by the dedicated [Protocol mappers](#_protocol-mappers). There is a [Built-in client scope **roles**](#_client_scopes_protocol) defined in the realm. Since it is a [Realm default client scope](#proc_updating_client_scopes_server_administration_guide), it is defined by default as a [Default client scope](#_client_scopes_linking) for every realm client. You can see this client scope in the admin console by looking at the tab **Client scopes** and then looking for the **roles** client scope. This client scope contains these protocol mappers by default:

- The protocol mapper **realm roles** - This protocol mapper is used to add the realm roles to the token claim. By default, the configuration looks like this:

Realm roles mapper

![mapper oidc realm roles](./images/mapper-oidc-realm-roles.png)

- The protocol mapper **client roles** - This protocol mapper is used to add the client roles to the token claim. By default, the configuration looks like this:

Client roles mapper

![mapper oidc client roles](./images/mapper-oidc-client-roles.png)

- The protocol mapper **audience resolve** - This protocol mapper is used to fill the `aud` claim in the access token based on the applied client roles. The details about this mapper are in the [Audience resolve section](#_audience_resolve).

As you can see in the configuration of realm roles and client roles mappers, it is possible to configure:

- If roles are added just to the access token or also to other tokens, like for example the ID token. By default, roles are added to the access token and to the introspection endpoint.
- What are the claims where the roles would be added. By default, the realm roles are added to the `realm_access` claim. So, for example, the claim in the JWT token containing 2 realm roles `role1` and `role2` will look similar to this:
  
  ```
  "realm_access": {
    "roles": [ "role1", "role2" ]
  }
  ```
  
  The client roles are added to the `resource_access` token claim by default. This claim will look like this in the token, which contains client roles `manage-account` and `manage-account-links` of client `account` and client role `target-client1-role` of the client `target-client1`:
  
  ```
  "resource_access": {
    "target-client1": {
      "roles": [ "target-client1-role" ]
    },
    "account": {
      "roles": [ "manage-account", "manage-account-links" ]
    }
  }
  ```

By adjusting the configuration option **Token claim name** of the role protocol mappers, it is possible to specify that these roles will be added to the token in the configured claim.

If you want to update the role claims just for one specific client (For example, client `foo` expects the realm roles in the claim `my-realm-roles` instead of the claim `realm_access`), then it is possible to remove the default client scope **roles** from your client and instead configure the realm/client protocol mapper in the [dedicated client scope](#_client_scopes_dedicated) of your client.

##### [](#example)Example

The [Audience documentation](#_audience_resolve) contains a more detailed example, which covers some details about the role mappings and about the audience (Claim `aud`) added to the token. Also, it can be useful to try the [Client scopes evaluation](#_client_scopes_evaluate) to see what are the effective scopes, protocol mappers and role scope mappings used when issuing the token for the particular client and how the JWT tokens would look like for the particular combination of user, client, and applied client scopes.

#### [](#audience-support)Audience support

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/oidc/con-audience.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-audience.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Foidc%2Fcon-audience.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Typically, the environment where Keycloak is deployed consists of a set of *confidential* or *public* client applications that use Keycloak for authentication. These clients are *frontend clients*, which may directly redirect user to Keycloak to request browser authentication. The particular client would then receive set of tokens after successful authentication.

*Services* (*Resource Servers* in the [OAuth 2 specification](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-mtls-08#section-4.2)) are also available that serve requests from client applications and provide resources to these applications. These services require an *Access token* (Bearer token) to be sent to them from *frontend application* or from other service to authenticate a request.

The care must be taken to make sure that access tokens have limited privileges and the particular access token cannot be misused by the service to access other third-party services. In the environment where trust among services is low, you may encounter this example scenario:

1. A frontend client application `frontend-client` requires authentication against Keycloak.
2. Keycloak authenticates a user.
3. Keycloak issues a token to the application `frontend-client`.
4. The `frontend-client` application uses the token to invoke a service `service1`.
5. The `service1` service returns the response to the application. But assume that this service will try to misuse the token and keep it for the further use.
6. The `service1` then invokes another service `service2` using the applications token, which was previously sent to it. The `service2` does not check that token was not supposed to be used to invoke it and it will serve the request and return successful response. This results in broken security as the `service1` misused the token to access other services on behalf of the client application `frontend-client`.

This scenario is unlikely in environments with a high level of trust between services but not in environments where trust is low.

To prevent any misuse of the access token, the access token can contain the claim `aud`, which represents the audience. The claim `aud` should typically represent client ids of all services where the token is supposed to be used. In the environments with low trust among services, it is recommended to:

- Limit the audience on the token to make sure that access tokens contain just limited amount of audiences.
- Configure your services to verify the audience on the token.

To prevent `service1` from the example above to misuse the token, the secure variant of the flow may instead look like this:

1. A frontend application `frontend-client` authenticates against Keycloak.
2. Keycloak authenticates a user.
3. Keycloak issues a token to the `frontend-client` application. The `frontend-client` knows that it will need to invoke `service1` so it places `scope=service1-scope` in the authentication request sent to Keycloak. The scope `service1-scope` is a [Client scope](#_client_scopes), which may need to be created by administrator. In the [sections below](#_audience_setup) there are some options how to setup such a client scope. The token claim will look like:
   
   ```
   "aud": "service1"
   ```
   
   This declares that the client can use this access token to invoke the `service1`.
4. The `frontend-client` application uses the token to invoke a service `service1`.
5. The `service1` serves the request to the client application `frontend-application`. But assume that this service will try to misuse the token and keep it for the further use.
6. The `service1` will then try to invoke a `service2` with the token. Invocation is not successful because the `service2` service checks the audience on the token and find that its audience is only for the `service1`. Hence `service2` will reject the request and will return an error to `service1`. This behavior is expected and security is not broken.

##### [](#ability-for-the-service-to-call-another-service)Ability for the service to call another service

In some environments, it may be desired that the `service1` may have to retrieve additional data from a `service2` to return data to the original client application `frontend-client`. In order to make this possible to work, there are few possibilities:

- Make sure that initial access token issued to `frontend-client` will contain both `service1` and `service2` as audiences. Assuming that there are proper client scopes set, the `frontend-client` can possibly use the `scope=service1-scope service2-scope` as a value of the `scope` parameter. The issued token would then contain the `aud` claim like:
  
  ```
  "aud": [ "service1", "service2" ]
  ```
  
  Such access token can be used to invoke both `service1` or `service2`. Hence `service1` will be able to successfully use such token to invoke `service2` to retrieve additional data.
- The previous approach with both services in the token audience allows that `service1` is allowed to invoke `service2`. However it means that `frontend-client` can also directly use his access token to invoke `service2`. This may not be desired in some cases. You may want `service1` to be able to invoke `service2`, but at the same time, you do not want `frontend-client` to be able to directly invoke `service2`. The solution to such scenario might be the use of the [Token exchange](https://www.keycloak.org/securing-apps/token-exchange). In that case, the initial token would still have only `service1` as an audience. However once the token is sent to `service1`, the `service1` may send Token exchange request to exchange the token for another token, which would have `service2` as an audience. Please see the [Token exchange Documentation](https://www.keycloak.org/securing-apps/token-exchange) for the details on how to use it.

##### [](#_audience_setup)Setup

When setting up audience checking:

- Ensure that services are configured to check audience on the access token sent to them. This may be done in a way specific to your client OIDC adapter, which you are using to secure your OIDC client application.
- Ensure that access tokens issued by Keycloak contain all necessary audiences.
  
  Audiences can be added to the token by two ways:
  
  - Using the client roles as described in the [Audience resolve section](#_audience_resolve).
  - Hardcoded audience as described in the [Hardcoded audience section](#_audience_hardcoded).

##### [](#_audience_resolve)Automatically add audience based on client roles

An *Audience Resolve* protocol mapper is defined in the default client scope *roles*. The mapper checks for clients that have at least one client role available for the current token. The client ID of each such client is then added as an audience, which is useful if your service clients rely on client roles. Service client could be usually a client without any flows enabled, which may not have any tokens issued directly to itself. It represents an OAuth 2 *Resource Server*.

The [Token role mappings section](#_oidc_token_role_mappings) contains the details about how are client roles added into the token. Please also see the example below.

###### [](#example-token-role-mappings-and-audience-claim)Example - token role mappings and audience claim

Here are the example steps how to use the client roles to make `aud` claim added to the token:

1. Create a [OIDC client](#proc-creating-oidc-client_server_administration_guide) `service1`. It may be possible to disable **Standard flow** or any other flows for this client as it is a service client, which may never directly authenticate by itself. The possible exception might be **Standard Token Exchange** switch if needed as described above.
2. Go to **Roles** tab of that client and create client role `service1-role`.
3. Create user `john` in the same realm and assign him the client role `service1-role` of client `service1` created in the previous step. [This section](#proc-assigning-role-mappings_server_administration_guide) contains some details on how to do it.
4. Create client scope named `service1-scope`. It can be marked with **Include in token scope** as **ON**. See [this section](#_client_scopes) for the details on how to create and set new client scope.
5. Go to the tab **Scope** of the `service1-scope` and add the role `service1-role` of the client `service1` to the [Role scope mappings](#_role_scope_mappings) of this client scope
6. Create another client `frontend-client` in the realm.
7. Click to the tab **Client scopes** of this client and select the first dedicated client scope `frontend-client-dedicated` and then go to the tab **Scope** and disable **Full scope allowed** switch
8. Go back to the tab **Client scopes** of this client and click **Add client scope** and link the `service1-scope` as **Optional**. See [Client Scopes Linking section](#_client_scopes_linking) for more details.
9. Click the sub-tab **Evaluate** in the **Client scopes** as described in [this section](#_client_scopes_evaluate). When filling user `john` and the subtab **Generated access token**, it can be seen that there is not any `aud` claim as there are not any client roles in the generated example token. However when adding also the scope `service1-scope` to the **Scope** field, it can be seen that there is client role `service1-role` as it is in **Role scope mappings** of the `service1-scope` and also in the role mappings of the user `john`. Due to that the `aud` claim will also contain `service1`.

Audience resolve example

![audience resolving evaluate](./images/audience_resolving_evaluate.png)

If you want the `service1` audience to be always applied for the tokens issued to the `frontend-client` client (without using the parameter `scope=service1-scope`), it can be fine to instead do any of these:

- Assign the `service1-scope` as **Default** client scope rather than **Optional**
- Add the role scope mapping of the `service1-role` directly to the [Dedicated client scope](#_client_scopes_dedicated) of the client. In this case, you will not need the `service1-scope` at all.

Note that since this approach is based on client roles, it also requires that user himself (user `john` in the example above) is a member of some client role of the client `service1`. Otherwise if there are not any client roles assigned, the audience `service1` will not be included. If you want audience to be included regardless of client roles, see the [Hardcoded audience](#_audience_hardcoded) section instead.

The frontend client itself is not automatically added to the access token audience, therefore allowing easy differentiation between the access token and the ID token, since the access token will not contain the client for which the token is issued as an audience.

If you need the client itself as an audience, see the [hardcoded audience](#_audience_hardcoded) option. However, using the same client as both frontend and REST service is not recommended.

##### [](#_audience_hardcoded)Hardcoded audience

When your service relies on realm roles or does not rely on the roles in the token at all, it can be useful to use a hardcoded audience. A hardcoded audience is a protocol mapper, that will add the client ID of the specified service client as an audience to the token. You can use any custom value, for example a URL, if you want to use a different audience than the client ID.

You can add the protocol mapper directly to the frontend client. If the protocol mapper is added directly, the audience will always be added as well.

For more control over the protocol mapper, you can create the protocol mapper on the dedicated client scope, which will be called for example **service2**.

Here the example steps for the hardcoded audience

1. Create a client `service2`
2. Create a client scope `service2-scope`.
3. In the tab **Mappers** of that client scope, select **Configure a new mapper** and select **Audience**
4. Select **Included Client Audience** as a `service2` and save the mapper
   
   Audience protocol mapper
   
   ![audience mapper](./images/audience_mapper.png)
5. Link the newly created client scope with some client. For example it can be linked as **Optional** client scope to the client `frontend-client` created in the [previous example](#_audience_resolve).
6. You can optionally [Evaluate Client Scopes](#_client_scopes_evaluate) for the client where the client scope was linked (For example `frontend-client`) and generate an example access token. The audience `service2` will be added to the audience of the generated access token if `service2-scope` is included in the *scope* parameter, when you assigned it as an optional client scope.

In your confidential client application, ensure that the *scope* parameter is used. The value like *scope=service2-scope* must be included when you want to issue the token for accessing `service2`.

See in the [Keycloak JavaScript adapter](https://www.keycloak.org/securing-apps/javascript-adapter) section if your application uses the javascript adapter for how to send the *scope* parameter with the desired value.

If you prefer to not include `scope` parameter in your requests, you can instead link the `service2-scope` as a **Default** client scope or use the client dedicated scope where you configure this mapper. This is useful if you want to always apply the audience for all the authentication request of OIDC client `frontend-client`.

Both the *Audience* and *Audience Resolve* protocol mappers add the audiences to the access token only, by default. The ID Token typically contains only a single audience, the client ID for which the token was issued, a requirement of the OpenID Connect specification. However, the access token does not necessarily have the client ID, which was the token issued for, unless the *Audience* mapper added it.

##### [](#token-introspection-audience-validation)Token introspection audience validation

The OAuth2 token introspection endpoint validates that the authenticated client is present in the token’s audience (`aud`) claim before allowing introspection. This prevents clients from introspecting tokens that were not intended for them.

If you need to disable this validation during migration, you can enable a backwards compatibility option in **OpenID Connect Compatibility Modes** by enabling **Allow token introspection without audience check** on the client that performs the introspection. See [OpenID Connect provider configuration](https://www.keycloak.org/server/all-provider-config#_openid_connect) for the server-wide option.

### [](#_client-saml-configuration)Creating a SAML client

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/saml/proc-creating-saml-client.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fsaml%2Fproc-creating-saml-client.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fsaml%2Fproc-creating-saml-client.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak supports [SAML 2.0](#_saml) for registered applications. POST and Redirect bindings are supported. You can choose to require client signature validation. You can have the server sign and/or encrypt responses as well.

Procedure

1. Click **Clients** in the menu.
2. Click **Create client** to go to the **Create client** page.
3. Set **Client type** to **SAML**.
   
   Create client
   
   ![add client saml](./images/add-client-saml.png)
4. Enter the **Client ID** of the client. This is often a URL and is the expected **issuer** value in SAML requests sent by the application.
5. Click **Save**. This action creates the client and brings you to the **Settings** tab.

The following sections describe each setting on this tab.

#### [](#settings-tab)Settings tab

The **Settings** tab includes many options to configure this client.

Client settings

![client settings saml](./images/client-settings-saml.png)

##### [](#general-settings-2)General settings

**Client ID**

The alphanumeric ID string that is used in OIDC requests and in the Keycloak database to identify the client. This value must match the issuer value sent with AuthNRequests. Keycloak pulls the issuer from the Authn SAML request and match it to a client by this value.

**Name**

The name for the client in a Keycloak UI screen. To localize the name, set up a replacement string value. For example, a string value such as ${myapp}. See the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more information.

**Description**

The description of the client. This setting can also be localized.

**Always Display in Console**

Always list this client in the Account Console even if this user does not have an active session.

##### [](#access-settings-2)Access Settings

**Root URL**

When Keycloak uses a configured relative URL, this value is prepended to the URL.

**Home URL**

If Keycloak needs to link to a client, this URL is used.

**Valid Redirect URIs**

Enter a URL pattern and click the + sign to add. Click the - sign to remove. Click **Save** to save these changes. Wildcards values are allowed only at the end of a URL. For example, [http://host.com/\*$$](http://host.com/*$$). This field is used when the exact SAML endpoints are not registered and Keycloak pulls the Assertion Consumer URL from a request.

**IDP-Initiated SSO URL name**

URL fragment name to reference client when you want to do IDP Initiated SSO. Leaving this empty will disable IDP Initiated SSO. The URL you will reference from your browser will be: *server-root*/realms/{realm-name}/protocol/saml/clients/{client-url-name}

**IDP Initiated SSO Relay State**

Relay state you want to send with SAML request when you want to do IDP Initiated SSO.

**Master SAML Processing URL**

This URL is used for all SAML requests and the response is directed to the SP. It is used as the Assertion Consumer Service URL and the Single Logout Service URL.

If login requests contain the Assertion Consumer Service URL then those login requests will take precedence. This URL must be validated by a registered Valid Redirect URI pattern.

##### [](#saml-capabilities)SAML capabilities

**Name ID Format**

The Name ID Format for the subject. This format is used if no name ID policy is specified in a request, or if the Force Name ID Format attribute is set to ON.

**Force Name ID Format**

If a request has a name ID policy, ignore it and use the value configured in the Admin Console under **Name ID Format**.

**Force POST Binding**

By default, Keycloak responds using the initial SAML binding of the original request. By enabling **Force POST Binding**, Keycloak responds using the SAML POST binding even if the original request used the redirect binding.

**Force artifact binding**

If enabled, response messages are returned to the client through the SAML ARTIFACT binding system.

**Include AuthnStatement**

SAML login responses may specify the authentication method used, such as password, as well as timestamps of the login and the session expiration. **Include AuthnStatement** is enabled by default, so that the **AuthnStatement** element will be included in login responses. Setting this to OFF prevents clients from determining the maximum session length, which can create client sessions that do not expire.

**Include OneTimeUse Condition**

If enable, a OneTimeUse Condition is included in login responses.

**Optimize REDIRECT signing key lookup**

When set to ON, the SAML protocol messages include the Keycloak native extension. This extension contains a hint with the signing key ID. The SP uses the extension for signature validation instead of attempting to validate the signature using keys.

This option applies to REDIRECT bindings where the signature is transferred in query parameters and this information is not found in the signature information. This is contrary to POST binding messages where key ID is always included in document signature.

This option is used when Keycloak server and adapter provide the IDP and SP. This option is only relevant when **Sign Documents** is set to ON.

**Allow ECP Flow**

If true, this application is allowed to use SAML ECP profile for authentication.

##### [](#signature-and-encryption)Signature and Encryption

**Sign Documents**

When set to ON, Keycloak signs the document using the realms private key.

**Sign Assertions**

The assertion is signed and embedded in the SAML XML Auth response.

**Signature Algorithm**

The algorithm used in signing SAML documents. Note that `SHA1` based algorithms are deprecated and may be removed in a future release. We recommend the use of some more secure algorithm instead of `*_SHA1`. Also, with `*_SHA1` algorithms, verifying signatures do not work if the SAML client runs on Java 17 or higher.

**SAML Signature Key Name**

Signed SAML documents sent using POST binding contain the identification of the signing key in the **KeyName** element. This action can be controlled by the **SAML Signature Key Name** option. This option controls the contents of the **Keyname**.

- **KEY\_ID** The **KeyName** contains the key ID. This option is the default option.
- **CERT\_SUBJECT** The **KeyName** contains the subject from the certificate corresponding to the realm key. This option is expected by Microsoft Active Directory Federation Services.
- **NONE** The **KeyName** hint is completely omitted from the SAML message.

**Canonicalization Method**

The canonicalization method for XML signatures.

**Metadata descriptor URL**

External URL where the client publishes the `SPSSODescriptor` metadata. This URL is used to download the client certificates when the **Use metadata descriptor URL** is enabled.

**Use metadata descriptor URL**

When **ON**, the certificates to validate signatures (**Client signature required** option is enabled in the **Keys** tab) and encrypt assertions (**Encrypt assertions** in the same tab) are automatically downloaded from the `Metadata descriptor URL` and cached in Keycloak. If a specific certificate is requested to validate a signature (usually in `POST` binding) and it is not in the cache, certificates are automatically refreshed from the URL. If all certificates are requested to validate the signature (`REDIRECT` binding) or any key is requested to encrypt, the refresh is only done after a max cache time. This maximum time can be specified in the descriptor itself, `cacheDuration` or `validUntil` attributes, or the cache provider defines one. See [public-key-storage](https://www.keycloak.org/server/all-provider-config) spi in the all provider config guide for more information about how the cache works.

When the option is **OFF**, the key should be generated or imported when activating the respective switch in the **Keys** tab.

**Encryption algorithm**

Encryption algorithm used for the client. Default value is `AES_256_GCM` when not defined.

**Key transport algorithm**

Key transport algorithm used for the client to encrypt the secret key used for encryption. Default value is `RSA-OAEP-11` when not defined.

**Digest method for RSA-OAEP**

Digest method to use when RSA-OAEP is selected as the key transport algorithm. Only available if **Key transport algorithm** is set to any RSA-OAEP algorithm. Default value is `SHA-256` when not defined.

**Mask generation function**

Mask generation function to use when `RSA-OAEP-11` is selected as the key transport algorithm. Only available if **Key transport algorithm** is set to `RSA-OAEP-11` algorithm. Default value is `mgf1sha256` when no defined.

The encryption options are only available if the **Encrypt Assertions** option is enabled in the **Keys** tab. For more information about SAML/XML encryption, see the [XML Encryption Syntax and Processing](https://www.w3.org/TR/xmlenc-core1/) specification.

##### [](#login-settings-2)Login settings

**Login theme**

A theme to use for login, OTP, grant registration, and forgotten password pages.

**Consent required**

If enabled, users have to consent to client access.

For client-side clients that perform browser logins. As it is not possible to ensure that secrets can be kept safe with client-side clients, it is important to restrict access by configuring correct redirect URIs.

**Display client on screen**

This switch applies if **Consent Required** is **Off**.

- *Off*
  
  The consent screen will contain only the consents corresponding to configured client scopes.
- *On*
  
  There will be also one item on the consent screen about this client itself.

**Client consent screen text**

Applies if **Consent required** and **Display client on screen** are enabled. Contains the text that will be on the consent screen about permissions for this client.

##### [](#logout-settings-2)Logout settings

**Front channel logout**

If **Front Channel Logout** is enabled, the application requires a browser redirect to perform a logout. For example, the application may require a cookie to be reset which could only be done via a redirect. If **Front Channel Logout** is disabled, Keycloak invokes a background SAML request to log out of the application.

#### [](#keys-tab)Keys tab

**Client Signature Required**

If **Client Signature Required** is enabled, documents coming from a client are expected to be signed. Keycloak will validate this signature.

If the option **Use metadata descriptor URL** is enabled in the **Signature and Encryption** section of the **Settings** tab, the public keys used to validate signature are automatically downloaded and cached by Keycloak. If that option is disabled, you need to import or generate the key when **Client Signature Required** is activated.

**Encrypt Assertions**

Encrypts the assertions in SAML documents with the specified client public key. Default algorithms used for encryption are configured with security in mind. If you need a different configuration, the encryption details can be modified in the **Settings** tab, section **Signature and Encryption**. The encryption options are only visible when this **Encrypt Assertions** option is enabled.

The key used to encrypt the assertions is controlled in the same way as in the case of **Client Signature Required**. If **Use metadata descriptor URL** is enabled, the key is doenloaded and cached by Keycloak. If that option is disabled, you need to import or generate the key when activating the **Encrypt Assertions** option.

#### [](#advanced-tab-2)Advanced tab

This tab has many fields for specific situations. Some fields are covered in other topics. For details on other fields, click the question mark icon.

##### [](#fine-grain-saml-endpoint-configuration)Fine Grain SAML Endpoint Configuration

**Logo URL**

URL that references a logo for the Client application.

**Policy URL**

URL that the Relying Party Client provides to the End-User to read about how the profile data will be used.

**Terms of Service URL**

URL that the Relying Party Client provides to the End-User to read about the Relying Party’s terms of service.

**Assertion Consumer Service POST Binding URL**

POST Binding URL for the Assertion Consumer Service.

**Assertion Consumer Service Redirect Binding URL**

Redirect Binding URL for the Assertion Consumer Service.

**Logout Service POST Binding URL**

POST Binding URL for the Logout Service.

**Logout Service Redirect Binding URL**

Redirect Binding URL for the Logout Service.

**Logout Service Artifact Binding URL**

*Artifact* Binding URL for the Logout Service. When set together with the `Force Artifact Binding` option, *Artifact* binding is forced for both login and logout flows. *Artifact* binding is not used for logout unless this property is set.

**Logout Service SOAP Binding URL**

Redirect Binding URL for the Logout Service. Only applicable if **back channel logout** is used.

**Artifact Binding URL**

URL to send the HTTP artifact messages to.

**Artifact Resolution Service**

URL of the client SOAP endpoint where to send the `ArtifactResolve` messages to.

##### [](#advanced-settings)Advanced settings

**Assertion Lifespan**

Specific client lifespan set in the SAML assertion conditions. After that time the assertion will be invalid. If not specified the realm **Access Token Lifespan** is used. The `SessionNotOnOrAfter` attribute is not modified and continue using the **SSO Session Max** time defined at realm level.

**ACR to LoA Mapping**

Define which ACR (Authentication Context Class Reference) value is mapped to which LoA (Level of Authentication). The ACR for SAML is an URI, whereas the LoA must be numeric. This mapping overrides the [ACR to Level of Authentication (LoA) Mapping](#_mapping-acr-to-loa-realm) defined at realm level. Only present if [Step-up authentication for SAML](#_step-up-authentication-saml) feature is enabled.

**Minimum ACR Value**

Minimum ACR to be enforced by Keycloak. If the resulting authentication context for the request is as strong as this ACR the request is valid, otherwise Keycloak returns the `NoAuthnContext` status error. Only present if [Step-up authentication for SAML](#_step-up-authentication-saml) feature is enabled.

#### [](#idp-initiated-login)IDP Initiated login

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/saml/idp-initiated-login.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fsaml%2Fidp-initiated-login.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fsaml%2Fidp-initiated-login.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

IDP Initiated Login is a feature that allows you to set up an endpoint on the Keycloak server that will log you into a specific application/client. In the **Settings** tab for your client, you need to specify the **IDP Initiated SSO URL Name**. This is a simple string with no whitespace in it. After this you can reference your client at the following URL: `root/realms/{realm-name}/protocol/saml/clients/{url-name}`

The IDP initiated login implementation prefers *POST* over *REDIRECT* binding (check [saml bindings](#_saml) for more information). Therefore the final binding and SP URL are selected in the following way:

1. If the specific **Assertion Consumer Service POST Binding URL** is defined (inside **Fine Grain SAML Endpoint Configuration** section of the client settings) *POST* binding is used through that URL.
2. If the general **Master SAML Processing URL** is specified then *POST* binding is used again throughout this general URL.
3. As the last resort, if the **Assertion Consumer Service Redirect Binding URL** is configured (inside **Fine Grain SAML Endpoint Configuration**) *REDIRECT* binding is used with this URL.

If your client requires a special relay state, you can also configure this on the **Settings** tab in the **IDP Initiated SSO Relay State** field. Alternatively, browsers can specify the relay state in a **RelayState** query parameter, i.e. `root/realms/{realm-name}/protocol/saml/clients/{url-name}?RelayState=thestate`.

When using [identity brokering](#_identity_broker), it is possible to set up an IDP Initiated Login for a client from an external IDP. The actual client is set up for IDP Initiated Login at broker IDP as described above. The external IDP has to set up the client for application IDP Initiated Login that will point to a special URL pointing to the broker and representing IDP Initiated Login endpoint for a selected client at the brokering IDP. This means that in client settings at the external IDP:

- **IDP Initiated SSO URL Name** is set to a name that will be published as IDP Initiated Login initial point,
- **Assertion Consumer Service POST Binding URL** in the **Fine Grain SAML Endpoint Configuration** section has to be set to the following URL: `broker-root/realms/{broker-realm}/broker/{idp-name}/endpoint/clients/{client-id}`, where:
  
  - *broker-root* is base broker URL
  - *broker-realm* is name of the realm at broker where external IDP is declared
  - *idp-name* is name of the external IDP at broker
  - *client-id* is the value of **IDP Initiated SSO URL Name** attribute of the SAML client defined at broker. It is this client, which will be made available for IDP Initiated Login from the external IDP.

Please note that you can import basic client settings from the brokering IDP into client settings of the external IDP - just use [SP Descriptor](#_identity_broker_saml_sp_descriptor) available from the settings of the identity provider in the brokering IDP, and add `clients/client-id` to the endpoint URL.

#### [](#proc-using-an-entity-descriptors_server_administration_guide)Using an entity descriptor to create a client

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/saml/proc-using-an-entity-descriptor.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fsaml%2Fproc-using-an-entity-descriptor.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fsaml%2Fproc-using-an-entity-descriptor.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Instead of registering a SAML 2.0 client manually, you can import the client using a standard SAML Entity Descriptor XML file.

The Client page includes an **Import client** option.

Add client

![Import SAML client](./images/import-client-saml.png)

Procedure

1. Click **Browse**.
2. Load the file that contains the XML entity descriptor information.
3. Review the information to ensure everything is set up correctly.

Some SAML client adapters, such as *mod-auth-mellon*, need the XML Entity Descriptor for the IDP. You can find this descriptor by going to this URL:

```
root/realms/{realm-name}/protocol/saml/descriptor
```

where *realm* is the realm of your client.

### [](#con-client-links_server_administration_guide)Client links

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/con-client-links.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fcon-client-links.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fcon-client-links.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To link from one client to another, Keycloak provides a redirect endpoint: `/realms/realm_name/clients/{client-id}/redirect`.

If a client accesses this endpoint using a `HTTP GET` request, Keycloak returns the configured base URL for the provided Client and Realm in the form of an `HTTP 307` (Temporary Redirect) in the response’s `Location` header. As a result of this, a client needs only to know the Realm name and the Client ID to link to them. This indirection avoids hard-coding client base URLs.

As an example, given the realm `master` and the client-id `account`:

```
http://host:port/realms/master/clients/account/redirect
```

This URL temporarily redirects to: [http://host:port/realms/master/account](http://host:port/realms/master/account)

### [](#_protocol-mappers)OIDC token and SAML assertion mappings

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/con-protocol-mappers.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fcon-protocol-mappers.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fcon-protocol-mappers.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Applications receiving ID tokens, access tokens, or SAML assertions may require different roles and user metadata.

You can use Keycloak to:

- Hardcode roles, claims and custom attributes.
- Pull user metadata into a token or assertion.
- Rename roles.

You perform these actions in the **Mappers** tab in the Admin Console.

Mappers tab

![mappers oidc](./images/mappers-oidc.png)

New clients do not have built-in mappers, but they can inherit some mappers from client scopes. See the [client scopes section](#_client_scopes) for more details.

Protocol mappers map items (such as an email address, for example) to a specific claim in the identity and access token. The function of a mapper should be self-explanatory from its name. You add pre-configured mappers by clicking **Add Builtin**.

Each mapper has a set of common settings. Additional settings are available, depending on the mapper type. Click **Edit** next to a mapper to access the configuration screen to adjust these settings.

Mapper config

![mapper config](./images/mapper-config.png)

Details on each option can be viewed by hovering over its tooltip.

You can use most OIDC mappers to control where the claim gets placed. You opt to include or exclude the claim from the *id* and *access* tokens by adjusting the **Add to ID token** and **Add to access token** switches.

You can add mapper types as follows:

Procedure

1. Go to the **Mappers** tab.
2. Click **Configure a new mapper**.
   
   Add mapper
   
   ![add mapper](./images/add-mapper.png)
3. Select a **Mapper Type** from the list box.

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/proc-creating-mappers.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fproc-creating-mappers.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fproc-creating-mappers.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

#### [](#_protocol-mappers_priority)Priority order

Mapper implementations have *priority order*. *Priority order* is not the configuration property of the mapper. It is the property of the concrete implementation of the mapper.

Mappers are sorted by the order in the list of mappers. The changes in the token or assertion are applied in that order with the lowest applying first. Therefore, the implementations that are dependent on other implementations are processed in the necessary order.

For example, to compute the roles which will be included with a token:

1. Resolve audiences based on those roles.
2. Process a JavaScript script that uses the roles and audiences already available in the token.

#### [](#_protocol-mappers_oidc-user-session-note-mappers)OIDC user session note mappers

User session details are defined using mappers and are automatically included when you use or enable a feature on a client. Click **Add builtin** to include session details.

Impersonated user sessions provide the following details:

- **IMPERSONATOR\_ID**: The ID of an impersonating user.
- **IMPERSONATOR\_USERNAME**: The username of an impersonating user.

Service account sessions provide the following details:

- **clientId**: The client ID of the service account.
- **client\_id**: The client ID of the service account.
- **clientAddress**: The remote host IP of the service account’s authenticated device.
- **clientHost**: The remote host name of the service account’s authenticated device.

#### [](#script-mapper)Script mapper

Use the **Script Mapper** to map claims to tokens by running user-defined JavaScript code. For more details about deploying scripts to the server, see [JavaScript Providers](https://www.keycloak.org/docs/26.6.3/server_development/#_script_providers).

When scripts deploy, you should be able to select the deployed scripts from the list of available mappers.

#### [](#pairwise-subject-identifier-mapper)Pairwise subject identifier mapper

Subject claim *sub* is mapped by default by **Subject (sub)** protocol mapper in the default client scope **basic**.

To use a pairwise subject identifier by using a protocol mapper such as **Pairwise subject identifier**, you can remove the **Subject (sub)** protocol mapper from the **basic** client scope. However it is not strictly needed as the **Subject (sub)** protocol mapper is executed before the **Pairwise subject identifier** mapper and hence the pairwise value will override the value added by the Subject mapper. This is due to the [priority](#_protocol-mappers_priority) of the Subject mapper. So the only advantage of removing the built-in **Subject (sub)** mapper might be to save a little bit of performance by avoiding the use of the protocol mapper, which may not have any effect.

#### [](#_using_lightweight_access_token)Using lightweight access token

The access token in Keycloak contains sensitive information, including Personal Identifiable Information (PII). Therefore, if the resource server does not want to disclose this type of information to third party entities such as clients, Keycloak supports lightweight access tokens that remove PII from access tokens. Further, when the resource server acquires the PII removed from the access token, it can acquire the PII by sending the access token to Keycloak’s token introspection endpoint.

Information that cannot be removed from a lightweight access token

Protocol mappers can controls which information is put onto an access token and the lightweight access token use the protocol mappers. Therefore, the following information cannot be removed from the lightweight access.  
`exp`, `iat`, `jti`, `iss`, `typ`, `azp`, `sid`, `scope`, `cnf`

Using a lightweight access token in Keycloak

By applying `use-lightweight-access-token` executor of [client policies](#_client_policies) to a client, the client can receive a lightweight access token instead of an access token. The lightweight access token contains a claim controlled by a protocol mapper where its setting `Add to lightweight access token`(default OFF) is turned ON. Also, by turning ON its setting `Add to token introspection` of the protocol mapper, the client can obtain the claim by sending the access token to Keycloak’s token introspection endpoint.

Introspection endpoint

In some cases, it might be useful to trigger the token introspection endpoint with the HTTP header `Accept: application/jwt` instead of `Accept: application/json`, which can be useful especially for lightweight access tokens. See the details of **Token Introspection endpoint** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section.

UserInfo endpoint restriction

The UserInfo endpoint rejects lightweight access tokens by default, not for direct use with UserInfo. If you need user information from a lightweight access token, use one of these alternatives:

- Call the token introspection endpoint instead of UserInfo - the introspection endpoint is designed for lightweight tokens and returns the full claims
- Exchange the lightweight access token for a full access token using [Token Exchange](https://www.keycloak.org/securing-apps/token-exchange), then call UserInfo with the full token
  
  If you must use lightweight tokens with the UserInfo endpoint during migration, you can enable a backwards compatibility option in **OpenID Connect Compatibility Modes** by enabling **Allow UserInfo with lightweight access token** on the client that issues the lightweight tokens.

### [](#_client_installation)Generating client adapter config

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/proc-generating-client-adapter-config.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fproc-generating-client-adapter-config.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fproc-generating-client-adapter-config.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak can generate configuration files that you can use to install a client adapter in your application’s deployment environment. A number of adapter types are supported for OIDC and SAML.

1. Click on the *Action* menu and select the **Download adapter config** option
   
   ![client installation](./images/client-installation.png)
2. Select the **Format Option** you want configuration generated for.

All Keycloak client adapters for OIDC and SAML are supported. The mod-auth-mellon Apache HTTPD adapter for SAML is supported as well as standard SAML entity descriptor files.

### [](#_client_scopes)Client scopes

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/con-client-scopes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fcon-client-scopes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fcon-client-scopes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Use Keycloak to define a shared client configuration in an entity called a *client scope*. A *client scope* configures [protocol mappers](#_protocol-mappers) and [role scope mappings](#_role_scope_mappings) for multiple clients.

Client scopes also support the OAuth 2 **scope** parameter. Client applications use this parameter to request claims or roles in the access token, depending on the requirement of the application.

To create a client scope, follow these steps:

1. Click **Client Scopes** in the menu.
   
   Client scopes list
   
   ![client scopes list](./images/client-scopes-list.png)
2. Click **Create**.
3. Name your client scope.
4. Click **Save**.

A *client scope* has similar tabs to regular clients. You can define [protocol mappers](#_protocol-mappers) and [role scope mappings](#_role_scope_mappings). These mappings can be inherited by other clients and are configured to inherit from this client scope.

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/proc-creating-client-scopes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fproc-creating-client-scopes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fproc-creating-client-scopes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

#### [](#_client_scopes_protocol)Protocol

When you create a client scope, choose the **Protocol**. Clients linked in the same scope must have the same protocol.

Each realm has a set of pre-defined built-in client scopes in the menu.

- SAML protocol: The **role\_list**. This scope contains one protocol mapper for the roles list in the SAML assertion.
- OpenID Connect protocol: Several client scopes are available:
  
  - **roles**
    
    This scope is not defined in the OpenID Connect specification and is not added automatically to the **scope** claim in the access token. This scope has mappers, which are used to add the roles of the user to the access token and add audiences for clients that have at least one client role. These mappers are described in more detail in the [Token Role mappings section](#_oidc_token_role_mappings) and [Audience section](#_audience_resolve).
  - **web-origins**
    
    This scope is also not defined in the OpenID Connect specification and not added to the **scope** claiming the access token. This scope is used to add allowed web origins to the access token **allowed-origins** claim.
  - **microprofile-jwt**
    
    This scope handles claims defined in the [MicroProfile/JWT Auth Specification](https://github.com/eclipse/microprofile/wiki/JWT_Auth). This scope defines a user property mapper for the **upn** claim and a realm role mapper for the **groups** claim. These mappers can be changed so different properties can be used to create the MicroProfile/JWT specific claims.
  - **offline\_access**
    
    This scope is used in cases when clients need to obtain offline tokens. More details on offline tokens is available in the [Offline Access section](#_offline-access) and in the [OpenID Connect specification](https://openid.net/specs/openid-connect-core-1_0.html#OfflineAccess).
  - **profile**
  - **email**
  - **address**
  - **phone**

The client scopes **profile**, **email**, **address** and **phone** are defined in the [OpenID Connect specification](https://openid.net/specs/openid-connect-core-1_0.html#ScopeClaims). These scopes do not have any role scope mappings defined but they do have protocol mappers defined. These mappers correspond to the claims defined in the OpenID Connect specification.

For example, when you open the **phone** client scope and open the **Mappers** tab, you will see the protocol mappers which correspond to the claims defined in the specification for the scope **phone**.

Client scope mappers

![client scopes phone](./images/client-scopes-phone.png)

When the **phone** client scope is linked to a client, the client automatically inherits all the protocol mappers defined in the **phone** client scope. Access tokens issued for this client contain the phone number information about the user, assuming that the user has a defined phone number.

Built-in client scopes contain the protocol mappers as defined in the specification. You are free to edit client scopes and create, update, or remove any protocol mappers or role scope mappings.

#### [](#consent-related-settings)Consent related settings

Client scopes contain options related to the consent screen. Those options are useful if the linked client if **Consent Required** is enabled on the client.

Display On Consent Screen

If **Display On Consent Screen** is enabled, and the scope is added to a client that requires consent, the text specified in **Consent Screen Text** will be displayed on the consent screen. This text is shown when the user is authenticated and before the user is redirected from Keycloak to the client. If **Display On Consent Screen** is disabled, this client scope will not be displayed on the consent screen.

Consent Screen Text

The text displayed on the consent screen when this client scope is added to a client when consent required defaults to the name of client scope. The value for this text can be customised by specifying a substitution variable with **${var-name}** strings. The customised value is configured within the property files in your theme. See the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/) for more information on customisation.

#### [](#include-in-token-scope)Include in token scope

There is the **Include in token scope** switch on the client scope. If on, the name of this client scope will be added to the access token property scope, and to the Token Response and Token Introspection Endpoint response claim `scope`. If off, this client scope will be omitted from the token and from the Token Introspection Endpoint response. As mentioned above, some built-in client scopes have this switch disabled, which means that they are not included in the `scope` claim even if they are applied for the particular request.

#### [](#_client_scopes_linking)Link client scope with the client

Linking between a client scope and a client is configured in the **Client Scopes** tab of the client. Here is how it looks for the client application `myclient`:

Client scopes linking to client

![client scopes default](./images/client-scopes-default.png)

There are two ways of linking between the client scope and the client.

Default Client Scopes

This setting is applicable to the OpenID Connect and SAML clients. Default client scopes are applied when issuing OpenID Connect tokens or SAML assertions for a client. The client will inherit Protocol Mappers and Role Scope Mappings that are defined on the client scope. For the OpenID Connect Protocol, the Mappers and Role Scope Mappings are always applied, regardless of the value used for the scope parameter in the OpenID Connect authorization request.

Optional Client Scopes

This setting is applicable only for OpenID Connect clients. Optional client scopes are applied when issuing tokens for this client but only when requested by the **scope** parameter in the OpenID Connect authorization request.

##### [](#example-2)Example

For this example, assume the client has **profile** and **email** linked as default client scopes, and **phone** and **address** linked as optional client scopes. The client uses the value of the scope parameter when sending a request to the OpenID Connect authorization endpoint.

```
scope=openid phone
```

The scope parameter contains the string, with the scope values divided by spaces. The value **openid** is the meta-value used for all OpenID Connect requests. The token will contain mappers and role scope mappings from the default client scopes **profile** and **email** as well as **phone**, an optional client scope requested by the scope parameter.

##### [](#_client_scopes_dedicated)Dedicated client scope

There is a special client scope, which is linked to every client. It is a dedicated client scope, which is always shown as the first client scope when you click on the tab **Client scopes** of the particular client. For example, for client `myclient`, the client scope is shown as `myclient-dedicated`. This client scope represents the protocol mappers and role scope mappings, which are linked directly to the client itself.

It is not possible to unlink the dedicated client scope from a client. Also, it is not possible to link this dedicated client scope to a different client. In other words, the dedicated client scope is useful just for protocol mappers and role scope mappings, which are specific to a single client. In case you want to share the same protocol mapper configuration among multiple clients, it is usually useful to create a client scope in the realm tab **Client scopes** and then link this shared client scope to every client that should apply this shared configuration.

In the tab **Scope** of the dedicated client scope, you can define role scope mappings applicable to this client. You can also see the switch **Full scope allowed** in this tab. The details about this switch are described in [this section](#_role_scope_mappings) and in [this section](#_oidc_token_role_mappings).

In the admin REST API and in the internal Keycloak storage, the dedicated client scope does not exist as its protocol mappers and role scope mappings are internally linked to the client itself. The. dedicated client scope is in fact just an abstraction for the admin console UI.

#### [](#_client_scopes_evaluate)Evaluating Client Scopes

The **Mappers** tab contains the protocol mappers and the **Scope** tab contains the role scope mappings declared for this client. They do not contain the mappers and scope mappings inherited from client scopes. It is possible to see the effective protocol mappers (that is the protocol mappers defined on the client itself as well as inherited from the linked client scopes) and the effective role scope mappings used when generating a token for a client.

Procedure

1. Click the **Client Scopes** tab for the client.
2. Open the sub-tab **Evaluate**.
3. Select the optional client scopes that you want to apply.

This will also show you the value of the **scope** parameter. This parameter needs to be sent from the application to the Keycloak OpenID Connect authorization endpoint.

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/proc-evaluating-client-scopes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fproc-evaluating-client-scopes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fproc-evaluating-client-scopes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

If your application uses the [Keycloak JavaScript adapter](https://www.keycloak.org/securing-apps/javascript-adapter), see its section to learn how to send the **scope** parameter with the desired value.

You can also simulate how the access token, ID token, or UserInfo response issued to this client looks for a particular selected user and for a specific value of the `audience` parameter. Note that the `audience` parameter is currently only supported for the token exchange grant. It is recommended to leave it empty when simulating any other grant.

Evaluating client scopes

![client scopes evaluate](./images/client-scopes-evaluate.png)

All examples are generated for the particular user and issued for the particular client, with the specified value of the **scope** parameter. The examples include all of the claims and role mappings used.

#### [](#client-scopes-permissions)Client scopes permissions

When issuing tokens to a user, the client scope applies only if the user is permitted to use it.

When a client scope does not have any role scope mappings defined, each user is permitted to use this client scope. However, when a client scope has role scope mappings defined, the user must be a member of at least one of the roles. There must be an intersection between the user roles and the roles of the client scope. Composite roles are factored into evaluating this intersection.

If a user is not permitted to use the client scope, no protocol mappers or role scope mappings will be used when generating tokens. The client scope will not appear in the *scope* value in the token.

#### [](#realm-default-client-scopes)Realm default client scopes

Use **Realm Default Client Scopes** to define sets of client scopes that are automatically linked to newly created clients.

To see the realm default client scopes, click the **Client Scopes** tab on the left side of the admin console. In the **Assigned type** column, you can specify whether a particular client scope should be added as a **Default Client Scope** or an **Optional Client Scope** to newly created clients. See [this section](#_client_scopes_linking) for details on what *default* and *optional* client scopes are.

When a client is created, you can unlink the default client scopes, if needed. This is similar to removing [Default Roles](#_default_roles).

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/proc-updating-default-scopes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fproc-updating-default-scopes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fproc-updating-default-scopes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

#### [](#_downscoping)Downscoping

In OAuth/OIDC, **downscoping** is the process of exchanging an existing JWT access token for a new one with a more restricted set of permissions (scopes) and/or a narrower audience. In the OAuth 2.0 refresh token grant type, the [RFC 6749](https://datatracker.ietf.org/doc/html/rfc6749#section-6) itself restricts the requested scope, saying that it must not include any scope not originally granted by the resource owner. So, this **downscoping** concept is very common and very recommended for security reasons.

Keycloak provides a [client policy](#_client_policies) executor that ensures this **downscoping** idea for other grant types. The client executor is called `downscope-assertion-grant-enforcer` and, for the moment, applies for the [Standard token exchange](https://www.keycloak.org/securing-apps/token-exchange#_standard-token-exchange). When this client executor is enforced, the token exchange is only allowed for the scopes that are already present in the initial JWT (`subject_token` parameter). An error is returned if any other extra scope is requested, no matter if the client configuration permits this scope as optional or default. Default scopes that are configured as **include in token scope** set to **false** (for example `basic` or `acr` in the default configuration) are the only exception. Those scopes are invisible for the requester and are considered compulsory for any grant type. Once this executor is applied for the client, **downscoping** is the only option when exchanging an access token, no additional scopes will ever be granted.

#### [](#scopes-explained)Scopes explained

The term *scope* has multiple meanings within Keycloak and across the OAuth/OIDC specifications. Below is a clarification of the different *scopes* used in Keycloak:

Client scope

Client scopes are entities in Keycloak that are configured at the realm level and can be linked to clients. Client scopes are referenced by their name when a request is sent to the Keycloak authorization endpoint with a corresponding value of the **scope** parameter. See the [client scopes linking](#_client_scopes_linking) section for more details.

Role scope mapping

This is available under the **Scope** tab of a client or client scope. Use **Role scope mapping** to limit the roles that can be used in the access tokens. See the [Role Scope Mappings section](#_role_scope_mappings) for more details.

Authorization scopes

The **Authorization Scope** covers the actions that can be performed in the application. See the [Authorization Services Guide](https://www.keycloak.org/docs/26.6.3/authorization_services/) for more details.

### [](#_client_policies)Client Policies

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/clients/client-policies.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fclients%2Fclient-policies.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fclients%2Fclient-policies.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To make it easy to secure client applications, it is beneficial to realize the following points in a unified way.

- Setting policies on what configuration a client can have
- Validation of client configurations
- Conformance to a required security standards and profiles such as Financial-grade API (FAPI) and OAuth 2.1

To realize these points in a unified way, *Client Policies* concept is introduced.

#### [](#use-cases)Use-cases

Client Policies realize the following points mentioned as follows.

Setting policies on what configuration a client can have

Configuration settings on the client can be enforced by client policies during client creation/update, but also during OpenID Connect requests to Keycloak server, which are related to particular client. Keycloak supports similar thing also through the **Client Registration Policies** described in the **Client registration service** in the [Securing applications and Services guide](https://www.keycloak.org/guides#securing-apps). However, Client Registration Policies can only cover OIDC Dynamic Client Registration. Client Policies cover not only what Client Registration Policies can do, but other client registration and configuration ways. The current plans are for Client Registration to be replaced by Client Policies.

Validation of client configurations

Keycloak supports validation whether the client follows settings like Proof Key for Code Exchange, Request Object Signing Algorithm, Holder-of-Key Token, and so on some endpoints like Authorization Endpoint, Token Endpoint, and so on. These can be specified by each setting item (on Admin Console, switch, pull-down menu and so on). To make the client application secure, the administrator needs to set many settings in the appropriate way, which makes it difficult for the administrator to secure the client application. Client Policies can do these validation of client configurations mentioned just above and they can also be used to autoconfigure some client configuration switches to meet the advanced security requirements. In the future, individual client configuration settings may be replaced by Client Policies directly performing required validations. Check the [dedicated section](#ssrf) to see a use case on how to protect the server by validating client uris.

Conformance to a required security standards and profiles such as FAPI and OAuth 2.1

The *Global client profiles* are client profiles pre-configured in Keycloak by default. They are pre-configured to be compliant with standard security profiles like **FAPI** and **OAuth 2.1** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section, which makes it easy for the administrator to secure their client application to be compliant with the particular security profile. At this moment, Keycloak has global profiles for the support of FAPI and OAuth 2.1 specifications. The administrator will just need to configure the client policies to specify which clients should be compliant with the FAPI and OAuth 2.1. The administrator can configure client profiles and client policies, so that Keycloak clients can be easily made compliant with various other security profiles like SPA, Native App, Open Banking and so on.

#### [](#protocol)Protocol

The client policy concept is independent of any specific protocol. Keycloak currently supports especially client profiles for the [OpenID Connect (OIDC) protocol](https://www.keycloak.org/docs/26.6.3/server_admin/#con-oidc_server_administration_guide), but there is also a client profile available for the [SAML protocol](https://www.keycloak.org/docs/26.6.3/server_admin/#_saml).

#### [](#architecture)Architecture

Client Policies consists of the four building blocks: Condition, Executor, Profile and Policy.

##### [](#condition)Condition

A condition determines to which client a policy is adopted and when it is adopted. Some conditions are checked at the time of client create/update when some other conditions are checked during client requests (OIDC Authorization request, Token endpoint request and so on). The condition checks whether one specified criteria is satisfied. For example, some condition checks whether the access type of the client is confidential.

The condition can not be used solely by itself. It can be used in a [policy](#_client_policy_policy) that is described afterwards.

A condition can be configurable the same as other configurable providers. What can be configured depends on each condition’s nature.

The following conditions are provided:

The way of creating/updating a client

- Dynamic Client Registration (Anonymous or Authenticated with Initial access token or Registration access token)
- Admin REST API (Admin Console and so on)

So for example when creating a client, a condition can be configured to evaluate to true when this client is created by OIDC Dynamic Client Registration without initial access token (Anonymous Dynamic Client Registration). So this condition can be used for example to ensure that all clients registered through OIDC Dynamic Client Registration are FAPI or OAuth 2.1 compliant.

Author of a client (Checked by presence to the particular role or group)

On OpenID Connect dynamic client registration, an author of a client is the end user who was authenticated to get an access token for generating a new client, not Service Account of the existing client that actually accesses the registration endpoint with the access token. On registration by Admin REST API, an author of a client is the end user like the administrator of the Keycloak.

Client Access Type (confidential, public, bearer-only)

For example when a client sends an authorization request, a policy is adopted if this client is confidential. Confidential client has enabled client authentication when public client has disabled client authentication. Bearer-only is a deprecated client type.

Client Scope

Evaluates to true if the client has a particular client scope (either as default or as an optional scope used in current request). This can be used for example to ensure that OIDC authorization requests with scope `fapi-example-scope` need to be FAPI compliant.

Client Role

Applies for clients with the client role of the specified name. Typically you can create a client role of specified name to requested clients and use it as a "marker role" to make sure that specified client policy will be applied for requested clients.

A use-case often exists for requiring the application of a particular client policy for the specified clients such as `my-client-1` and `my-client-2`. The best way to achieve this result is to use a **Client Role** condition in your policy and then a create client role of specified name to requested clients. This client role can be used as a "marker role" used solely for marking that particular client policy for particular clients.

Client Domain Name, Host or IP Address

Applied for specific domain names of client. Or for the cases when the administrator registers/updates client from particular Host or IP Address.

Client Attribute

Applies to clients with the client attribute of the specified name and value. If you specify multiple client attributes, they will be evaluated using AND conditions. If you want to evaluate using OR conditions, set this condition multiple times.

Any Client

This condition always evaluates to true. It can be used for example to ensure that all clients in the particular realm are FAPI compliant.

ACR Condition

Applied when an ACR value requested in the authentication request matches the value configured in the condition. For example, it can be used to select an authentication flow based on the requested ACR value. For more details, see the [related documentation](#_client-policy-auth-flow) and the [official OIDC specification](https://openid.net/specs/openid-connect-core-1_0.html#acrSemantics).

Grant Type

Evaluates to true when a specific grant type is used. For example, it can be used in combination with Client Scope to block a token exchange request when a specific client scope is requested.

Identity Provider Alias

Condition that checks the Identity Provider that is involved in the client request. A list of IdP alias can be configured. The condition evaluates to true if one of them is associated to the request. It only applies to operations in which an IdP is involved (for example JWT Authorization grant).

##### [](#executor)Executor

An executor specifies what action is executed on a client to which a policy is adopted. The executor executes one or several specified actions. For example, some executor checks whether the value of the parameter `redirect_uri` in the authorization request matches exactly with one of the pre-registered redirect URIs on Authorization Endpoint and rejects this request if not.

The executor can not be used solely by itself. It can be used in a [profile](#_client_policy_profile) that is described afterwards.

An executor can be configurable the same as other configurable providers. What can be configured depends on the nature of each executor.

An executor acts on various events. An executor implementation can ignore certain types of events (For example, executor for checking OIDC `request` object acts just on the OIDC authorization request). Events are:

- Creating a client (including creation through dynamic client registration)
- Updating a client
- Sending an authorization request
- Sending a token request
- Sending a token refresh request
- Sending a token revocation request
- Sending a token introspection request
- Sending a userinfo request
- Sending a logout request with a refresh token (note that logout with refresh token is proprietary Keycloak functionality unsupported by any specification. It is rather recommended to rely on the [official OIDC logout](#_oidc-logout)).

On each event, an executor can work in multiple phases. For example, on creating/updating a client, the executor can modify the client configuration by autoconfigure specific client settings. After that, the executor validates this configuration in validation phase.

One of several purposes for this executor is to realize the security requirements of client conformance profiles like FAPI and OAuth 2.1. To do so, the following executors are needed:

- Enforce secure [Client Authentication method](#_client-credentials) is used for the client
- Enforce [Holder-of-key tokens](#_mtls-client-certificate-bound-tokens) are used
- Enforce [Proof Key for Code Exchange (PKCE)](#_proof-key-for-code-exchange) is used
- Enforce secure signature algorithm for [Signed JWT client authentication (private-key-jwt)](#_client-credentials) is used
- Enforce HTTPS redirect URI and make sure that configured redirect URI does not contain wildcards
- Enforce OIDC `request` object satisfying high security level
- Enforce Response Type of OIDC Hybrid Flow including ID Token used as *detached signature* as described in the FAPI 1 specification, which means that ID Token returned from Authorization response won’t contain user profile data
- Enforce more secure `state` and `nonce` parameters treatment for preventing CSRF
- Enforce more secure signature algorithm when client registration
- Enforce `binding_message` parameter is used for CIBA requests
- Enforce [Client Secret Rotation](#_secret_rotation)
- Enforce Client Registration Access Token
- Enforce checking if a client is the one to which an intent was issued in a use case where an intent is issued before starting an authorization code flow to get an access token like UK OpenBanking
- Enforce prohibiting implicit and hybrid flow
- Enforce checking if a PAR request includes necessary parameters included by an authorization request
- Enforce [DPoP-binding tokens](#_dpop-bound-tokens) is used (available when `dpop` feature is enabled)
- Enforce [using lightweight access token](#_using_lightweight_access_token)
- Enforce that [refresh token rotation](#_refresh_token_rotation) is skipped and there is no refresh token returned from the refresh token response
- Enforce a valid redirect URI that the OAuth 2.1 specification requires
- Enforce SAML Redirect binding cannot be used or SAML requests and assertions are signed
- Enforce scopes granted in [Standard token exchange](https://www.keycloak.org/securing-apps/token-exchange#_standard-token-exchange) or in JWT Authorization Grant are restricted to the ones present in the initial `subject_token` or `assertion` JWT. This executor only allows downscoping of the presented assertion. An error is returned if any extra scope, not originally granted to the JWT, is requested.
- Enforce claims for assertion grants (`subject_token` in Token Exchange and `assertion` in JWT Authorization Grant). The executor enforces the presence and specific values of a claim in a JWT. It uses a Java regex so it is quite versatile.

Another available executor is the `auth-flow-enforce`, which can be used to enforce an authentication flow during an authentication request. For instance, it can be used to select a flow based on certain conditions, such as a specific scope or an ACR value. For more details, see the [related documentation](#_client-policy-auth-flow).

##### [](#_client_policy_profile)Profile

A profile consists of several executors, which can realize a security profile like FAPI and OAuth 2.1. Profile can be configured by the Admin REST API (Admin Console) together with its executors. Three *global profiles* exist and they are configured in Keycloak by default with pre-configured executors compliant with the FAPI 1 Baseline, FAPI 1 Advanced, FAPI CIBA, FAPI 2 and OAuth 2.1 specifications. More details exist in the **FAPI** and **OAuth 2.1** in the [securing apps](https://www.keycloak.org/guides#securing-apps) section.

##### [](#_client_policy_policy)Policy

A policy consists of several conditions and profiles. The policy can be adopted to clients satisfying all conditions of this policy. The policy refers several profiles and all executors of these profiles execute their task against the client that this policy is adopted to.

#### [](#configuration)Configuration

Policies, profiles, conditions, executors can be configured by Admin REST API, which means also the Admin Console. To do so, there is a tab *Realm* → *Realm Settings* → *Client Policies* , which means the administrator can have client policies per realm.

The *Global Client Profiles* are automatically available in each realm. However there are no client policies configured by default. This means that the administrator is always required to create any client policy if they want for example the clients of his realm to be FAPI compliant. Global profiles cannot be updated, but the administrator can easily use them as a template and create their own profile if they want to do some slight changes in the global profile configurations. There is JSON Editor available in the Admin Console, which simplifies the creation of new profile based on some global profile.

#### [](#backward-compatibility)Backward Compatibility

Client Policies can replace Client Registration Policies described in the **Client registration service** from [Securing applications Guides](https://www.keycloak.org/guides#securing-apps). However, Client Registration Policies also still co-exist. This means that for example during a Dynamic Client Registration request to create/update a client, both client policies and client registration policies are applied.

The current plans are for the Client Registration Policies feature to be removed and the existing client registration policies will be migrated into new client policies automatically.

#### [](#client-secret-rotation-example)Client Secret Rotation Example

See an example configuration for [client secret rotation](#_proc-secret-rotation).

#### [](#securing-client-uris)Securing Client URIs

For a client, it is possible to register many URIs for different purposes, such as the redirect URI for the Authorization Code Flow, JWKS URI to retrieve the client’s JWK Set, privacy URI, policy URI, logo URI, sector identifier URI, and others.

It is good practice to add as much protection as possible for these URIs to prevent malicious use. For example, you can enforce a URI to allow only specific hosts, to accept only HTTPS, or to match any desired pattern using regular expressions.

Keycloak, with client policies, provides specific executors to protect the various client URIs:

- **Secure Client URIs**: A generic executor that enforces all client URIs to only accept HTTPS and to not accept any wildcards.
- **Secure Redirect URIs Enforcer**: An executor that allows different levels of validations for redirect URIs, such as IPv4 and IPv6 loopback addresses, HTTP scheme, wildcards, only permitted domains, or ensuring OAuth 2.1 compliance.
- **Secure Client URIs Pattern Executor**: This executor is the most generic and flexible for validating client URLs. It allows you to select from a list only the URIs you want to validate, and then validates them providing a list of permitted pattern using regular expressions. Explore this [dedicated section](#ssrf) to see how this executor can be used to prevent an SSRF attack.

## [](#_oid4vci)Configuring Keycloak as a Verifiable Credential Issuer

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/oid4vci/vc-issuer-configuration.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Foid4vci%2Fvc-issuer-configuration.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Foid4vci%2Fvc-issuer-configuration.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

This is an experimental feature and should not be used in production. Backward compatibility is not guaranteed, and future updates may introduce breaking changes.

Keycloak provides experimental support for [OpenID for Verifiable Credential Issuance](https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0.html).

### [](#introduction)Introduction

This chapter provides step-by-step instructions for configuring Keycloak as a Verifiable Credential Issuer using the OpenID for Verifiable Credential Issuance (OID4VCI) protocol. It outlines the process for setting up a Keycloak instance to securely issue and manage Verifiable Credentials (VCs), supporting decentralized identity solutions.

### [](#what-are-verifiable-credentials-vcs)What are Verifiable Credentials (VCs)?

Verifiable Credentials (VCs) are cryptographically signed, tamper-evident data structures that represent claims about an entity, such as a person, organization, or device. They are foundational to decentralized identity systems, allowing secure and privacy-preserving identity verification without reliance on centralized authorities. VCs support advanced features like selective disclosure and zero-knowledge proofs, enhancing user privacy and security.

### [](#what-is-oid4vci)What is OID4VCI?

OpenID for Verifiable Credential Issuance (OID4VCI) is an extension of the OpenID Connect (OIDC) protocol. It defines a standardized, interoperable framework for credential issuers to deliver VCs to holders, who can then present them to verifiers. OID4VCI leverages Keycloak’s existing authentication and authorization capabilities to streamline VC issuance.

### [](#scope-of-this-chapter)Scope of This Chapter

This chapter covers the following technical configurations:

- Creating a dedicated realm for VC issuance.
- Setting up a test user for credential testing.
- Configuring custom cryptographic keys for signing and encrypting VCs.
- Defining realm attributes to specify VC metadata.
- Establishing client scopes and mappers to include user attributes in VCs.
- Registering a client to handle VC requests.
- Verifying the configuration using the issuer metadata endpoint.

### [](#prerequisites)Prerequisites

Ensure the following requirements are met before configuring Keycloak as a Verifiable Credential Issuer:

### [](#keycloak-instance)Keycloak Instance

A running Keycloak server with the OID4VCI feature enabled.

To enable the feature, add the following flag to the startup command:

```
--features=oid4vc-vci
```

Verify activation by checking the server logs for the `OID4VC_VCI` initialization message.

### [](#configuring-credential-issuance-in-keycloak)Configuring Credential Issuance in Keycloak

In Keycloak, Verifiable Credentials are managed through **ClientScopes**, with each ClientScope representing a single Verifiable Credential type. To enable the issuance of a credential, the corresponding ClientScope must be assigned to an OpenID Connect client - ideally as **optional**.

During the OAuth2 authorization process, the credential-specific scope can be requested by including the ClientScope’s name in the `scope` parameter of the authorization request. Once the user has successfully authenticated, the resulting Access Token **MUST** include the requested ClientScope in its `scope` claim. To ensure this, make sure the ClientScope option **Include in token scope** is enabled.

With this Access Token, the Verifiable Credential can be issued at the Credential Endpoint.

### [](#authentication)Authentication

An access token is required to authenticate API requests.

Refer to the following Keycloak documentation sections for detailed steps on:

- [Creating a Client](#proc-creating-oidc-client_server_administration_guide)
- [Obtaining an Access Token](#_oidc-auth-flows-direct)

### [](#configuration-steps)Configuration Steps

Follow these steps to configure Keycloak as a Verifiable Credential Issuer. Each section is detailed with procedures, explanations, and examples where applicable.

### [](#creating-a-realm)Creating a Realm

A realm in Keycloak is a logical container that manages users, clients, roles, and authentication flows. For Verifiable Credential (VC) issuance, create a dedicated realm to ensure isolation and maintain a clear separation of functionality.

For detailed instructions on creating a realm, refer to the Keycloak documentation: [Creating a Realm](#proc-creating-a-realm_server_administration_guide).

### [](#creating-a-user-account)Creating a User Account

A test user is required to simulate credential issuance and verify the setup.

For step-by-step instructions on creating a user, refer to the Keycloak documentation: [Creating a User](#assembly-managing-users_server_administration_guide).

Ensure that the user has a valid username, email, and password. If the password should not be reset upon first login, disable the "Temporary" toggle during password configuration.

### [](#key-management-configuration)Key Management Configuration

Keycloak uses cryptographic keys for signing and encrypting Verifiable Credentials (VCs). To ensure secure and standards-compliant issuance, configure **ECDSA (ES256) for signing**, **RSA (RS256) for signing**, and **RSA-OAEP for encryption** using a keystore.

For a detailed guide on configuring realm keys, refer to the Keycloak documentation: [Managing Realm Keys](#realm_keys).

#### [](#configuring-key-providers)Configuring Key Providers

To enable cryptographic operations for VC issuance:

- **ECDSA (ES256) Key**: Used for signing VCs with the ES256 algorithm.
- **RSA (RS256) Key**: Alternative signing mechanism using RS256.
- **RSA-OAEP Key**: Used for encrypting sensitive data in VCs.

Each key must be registered as a **java-keystore provider** within the **Realm Settings** &gt; **Keys** section, ensuring: - The keystore file is correctly specified and securely stored. - The appropriate algorithm (ES256, RS256, or RSA-OAEP) is selected. - The key is active, enabled, and configured with the correct usage (signing or encryption). - Priority values are set to define precedence among keys.

Ensure the keystore file is **securely stored** and accessible to the Keycloak server. Use **strong passwords** to protect both the keystore and the private keys.

### [](#registering-realm-attributes)Registering Realm Attributes

Realm attributes define metadata for Verifiable Credentials (VCs), such as **expiration times, supported formats, and scope definitions**. These attributes allow Keycloak to issue VCs with predefined settings.

Since the **Keycloak Admin Console does not support direct attribute creation**, use the **Keycloak Admin REST API** to configure these attributes.

#### [](#define-realm-attributes)Define Realm Attributes

Create a JSON file (e.g., `realm-attributes.json`) with the following content:

```
{
  "realm": "oid4vc-vci",
  "enabled": true,
  "attributes": {
    "preAuthorizedCodeLifespanS": 120
  }
}
```

#### [](#attribute-breakdown)Attribute Breakdown

The attributes section contains issuer-specific metadata: - **preAuthorizedCodeLifespanS** – Defines how long pre-authorized codes remain valid (in seconds). - **oid4vc.attestation.trusted\_keys** – JSON array of trusted JWK (JSON Web Key) objects for attestation proof validation. Each JWK must include a `kid` (key ID) field. These keys take precedence over realm session keys when there are conflicts. Useful for configuring additional trusted keys beyond the realm’s default keys. Format: JSON array of JWK objects, e.g., `[{"kid":"key1","kty":"EC",…​},{"kid":"key2","kty":"RSA",…​}]`. - **oid4vc.attestation.trusted\_key\_ids** – Comma-separated list of key IDs from the realm’s key providers to use for attestation proof validation. Keys are looked up by their `kid` regardless of enabled status, allowing the use of disabled keys that are not exposed in well-known endpoints. This attribute takes the highest priority when merging trusted keys. Format: comma-separated list of key IDs, e.g., `key-id-1,key-id-2,key-id-3`.

#### [](#import-realm-attributes)Import Realm Attributes

Use the following `curl` command to import the attributes into Keycloak:

```
curl -X PUT "https://localhost:8443/admin/realms/oid4vc-vci" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d @realm-attributes.json
```

- Replace `$ACCESS_TOKEN` with a valid **Keycloak Admin API access token**.
- **Avoid using `-k` in production**; instead, configure a **trusted TLS certificate**.

#### [](#time-claim-correlation-mitigation)Time-claim correlation mitigation

To reduce unintended correlation across multiple issuances or presentations, you can normalize time-related claims by either randomizing them within a time window or rounding them to a coarse time unit. This behavior is opt-in and controlled by the following realm attributes:

   Attribute Default Description

`oid4vci.time.claims.strategy`

`off`

Strategy to apply to time claims. Supported values: `off`, `randomize`, `round`.

`oid4vci.time.randomize.window.seconds`

`86400`

When strategy is `randomize`, subtract a random number of seconds between 0 and the value of this attribute from the original timestamp to mitigate correlation attacks.

`oid4vci.time.round.unit`

`SECOND`

When strategy is `round`, truncate timestamps to the selected unit boundary (UTC). Supported values: `SECOND`, `MINUTE`, `HOUR`, `DAY`.

How it is applied during issuance:

- For JWT-VC, the credential `issuanceDate` is normalized at issuance; the JWT `nbf` is derived from the normalized value. If a mapper sets a VC `expirationDate`, it is normalized and emitted as JWT `exp`.
- For SD-JWT VCs, time-related claims (`iat`, `nbf`, `exp`) are typically set using protocol mappers. Use the available OID4VC mappers, such as the Issued At Time Claim Mapper for `iat`, to populate these claims. Values are automatically normalized according to the realm strategy.

Examples:

```
# Round to start of day (UTC)
curl -X PUT "https://localhost:8443/admin/realms/oid4vc-vci" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
        "attributes": {
          "oid4vci.time.claims.strategy": "round",
          "oid4vci.time.round.unit": "DAY"
        }
      }'

# Randomize within the last 24 hours
curl -X PUT "https://localhost:8443/admin/realms/oid4vc-vci" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
        "attributes": {
          "oid4vci.time.claims.strategy": "randomize",
          "oid4vci.time.randomize.window.seconds": "86400"
        }
      }'
```

### [](#create-client-scopes-with-mappers)Create Client Scopes with Mappers

Client scopes define **which user attributes** are included in Verifiable Credentials (VCs). Therefore, they are considered the Verifiable Credential configuration itself. These scopes use **protocol mappers** to map specific claims into VCs and the protocol mappers will also contain the corresponding metadata for claims that is displayed at the Credential Issuer Metadata Endpoint.

You can create the ClientScopes using the Keycloak web Administration Console, but the web Administration Console does not yet support adding metadata configuration. For metadata configuration, you will need to use the Admin REST API.

#### [](#define-a-client-scope-with-a-mapper)Define a Client Scope with a Mapper

Create a JSON file (e.g., `client-scopes.json`) with the following content:

```
{
  "name": "vc-scope-mapping",
  "protocol": "oid4vc",
  "attributes": {
    "include.in.token.scope": "true",
    "vc.issuer_did": "did:web:vc.example.com",
    "vc.credential_configuration_id": "my-credential-configuration-id",
    "vc.credential_identifier": "my-credential-identifier",
    "vc.format": "jwt_vc",
    "vc.expiry_in_seconds": 31536000,
    "vc.verifiable_credential_type": "my-vct",
    "vc.supported_credential_types": "credential-type-1,credential-type-2",
    "vc.credential_contexts": "context-1,context-2",
    "vc.credential_signing_alg": "ES256",
    "vc.cryptographic_binding_methods_supported": "jwk",
    "vc.signing_key_id": "key-id-123456",
    "vc.display": "[{\"name\": \"IdentityCredential\", \"logo\": {\"uri\": \"https://university.example.edu/public/logo.png\", \"alt_text\": \"a square logo of a university\"}, \"locale\": \"en-US\", \"background_color\": \"#12107c\", \"text_color\": \"#FFFFFF\"}]",
    "vc.sd_jwt.number_of_decoys": "2",
    "vc.credential_build_config.sd_jwt.visible_claims": "iat,jti,nbf,exp,given_name",
    "vc.credential_build_config.hash_algorithm": "SHA-256",
    "vc.credential_build_config.token_jws_type": "JWS",
    "vc.include_in_metadata": "true"
  },
  "protocolMappers": [
    {
      "name": "academic_title-mapper-bsk",
      "protocol": "oid4vc",
      "protocolMapper": "oid4vc-static-claim-mapper",
      "config": {
        "claim.name": "academic_title",
        "staticValue": "N/A"
      }
    },
    {
      "name": "givenName",
      "protocol": "oid4vc",
      "protocolMapper": "oid4vc-user-attribute-mapper",
      "config": {
        "claim.name": "given_name",
        "userAttribute": "firstName",
        "vc.mandatory": "false",
        "vc.display": "[{\"name\": \"الاسم الشخصي\", \"locale\": \"ar-SA\"}, {\"name\": \"Vorname\", \"locale\": \"de-DE\"}, {\"name\": \"Given Name\", \"locale\": \"en-US\"}, {\"name\": \"Nombre\", \"locale\": \"es-ES\"}, {\"name\": \"نام\", \"locale\": \"fa-IR\"}, {\"name\": \"Etunimi\", \"locale\": \"fi-FI\"}, {\"name\": \"Prénom\", \"locale\": \"fr-FR\"}, {\"name\": \"पहचानी गई नाम\", \"locale\": \"hi-IN\"}, {\"name\": \"Nome\", \"locale\": \"it-IT\"}, {\"name\": \"名\", \"locale\": \"ja-JP\"}, {\"name\": \"Овог нэр\", \"locale\": \"mn-MN\"}, {\"name\": \"Voornaam\", \"locale\": \"nl-NL\"}, {\"name\": \"Nome Próprio\", \"locale\": \"pt-PT\"}, {\"name\": \"Förnamn\", \"locale\": \"sv-SE\"}, {\"name\": \"مسلمان نام\", \"locale\": \"ur-PK\"}]"
      }
    }
  ]
}
```

This is a **sample configuration**. You can define **additional protocol mappers** to support different claim mappings, such as:

- Dynamic attribute values instead of static ones.
- Mapping multiple attributes per credential type.
- Alternative supported credential types.

From the example above:

- It is important to set `include.in.token.scope=true`, see [Attribute table: include.in.token.scope](#include.in.token.scope).
- Most of the named attributes above are optional. See below: [Attribute Breakdown](#client-scope-attribute-breakdown).
- You can determine the appropriate `protocolMapper` names by first creating them through the Web Administration Console and then retrieving their definitions via the Admin REST API.

#### [](#client-scope-attribute-breakdown)Attribute Breakdown - ClientScope

   Property Required Description / Default

`name`

required

Name of the client scope.

`protocol`

required

Protocol used by the client scope. Use `oid4vc` for OpenID for Verifiable Credential Issuance, which is an OAuth2 extension (like `openid-connect`).

`include.in.token.scope`

required

[]()This value MUST be `true`. It ensures that the scope’s name is included in the `scope` claim of the issued Access Token.

`protocolMappers`

optional

Defines how claims are mapped into the credential and how metadata is exposed via the issuer’s metadata endpoint.

`vc.issuer_did`

optional

The Decentralized Identifier (DID) of the issuer.  
*Default*: `${name}`

`vc.credential_configuration_id`

optional

The credentials configuration ID.  
*Default*: `${name}+`

`vc.credential_identifier`

optional

The credentials identifier.  
*Default*: `${name}+`

`vc.format`

optional

Defines the VC format (e.g., `jwt_vc`).  
*Default*: `dc+sd-jwt`

`vc.verifiable_credential_type`

optional

The Verifiable Credential Type (VCT).  
*Default*: `${name}+`

`vc.supported_credential_types`

optional

The type values of the Verifiable Credential Type.  
*Default*: `${name}+`

`vc.credential_contexts`

optional

The context values of the Verifiable Credential Type.  
*Default*: `${name}+`

`vc.credential_signing_alg`

optional

Supported signature algorithm for this credential.  
*Default*: All asymmetric signing algorithms backed by realm keys.

`vc.cryptographic_binding_methods_supported`

optional

Supported cryptographic methods (if applicable).  
*Default*: `jwk`

`vc.signing_key_id`

optional

The ID of the key to sign this credential.  
*Default*: *none*

`vc.display`

optional

Display information shown in the user’s wallet about the issued credential.  
*Default*: *none*

`vc.sd_jwt.number_of_decoys`

optional

Used only with format `dc+sd-jwt`. Number of decoy hashes in the SD-JWT.  
*Default*: `10`

`vc.credential_build_config.sd_jwt.visible_claims`

optional

Used only with format `dc+sd-jwt`. Claims always disclosed in the SD-JWT body.  
*Default*: `id,iat,nbf,exp,jti`

`vc.credential_build_config.hash_algorithm`

optional

Hash algorithm used before signing the credential.  
*Default*: `SHA-256`

`vc.credential_build_config.token_jws_type`

optional

JWT type written into the `typ` header of the token.  
*Default*: `JWS`

`vc.expiry_in_s`

optional

Credential expiration time in seconds.  
*Default*: `31536000` (one year)

`vc.include_in_metadata`

optional

If this claim should be listed in the credentials metadata.  
*Default*: `true` but depends on the mapper-type. Claims like `jti`, `nbf`, `exp`, etc. are set to `false` by default.

`vc.key_attestations_required`

optional

Indicates whether issuing this credential requires a key attestation.  
*Default*: `false`.

`vc.key_attestations_required.key_storage`

optional

Comma separated list of accepted key-storage attack potential levels (see ISO 18045 levels, e.g. `iso_18045_high`).  
Only relevant if `vc.key_attestations_required` is present.  
*Default*: none

`vc.key_attestations_required.user_authentication`

optional

Comma separated list of accepted user-authentication attack potential levels (see ISO 18045 levels).  
Only relevant if `vc.key_attestations_required` is present.  
*Default*: none

#### [](#attribute-breakdown-protocolmappers)Attribute Breakdown - ProtocolMappers

- **name** – Mapper identifier.
- **protocol** – Must be `oid4vc` for Verifiable Credentials.
- **protocolMapper** – Specifies the claim mapping strategy (e.g., `oid4vc-static-claim-mapper`).
- **config**: contains the protocol-mappers specific attributes.

Most claims are dependent on the `protocolMapper`-value, but there are also commonly used claims available for all ProtocolMappers:

   Property Required Description / Default

`claim.name`

required

The name of the attribute that will be added into the Verifiable Credential.  
Just like with OIDC user attributes, you may use dots to create nested JSON objects.  
*Default*: *none*

`userAttribute`

required

The name of the users-attribute that will be used to map the value into the `claim.name` of the Verifiable Credential.  
*Default*: *none*

`vc.mandatory`

optional

If the credential must be issued with this claim.  
*Default*: `false`

`vc.display`

optional

Metadata information that is displayed at the credential-issuer metadata-endpoint.  
*Default*: *none*

#### [](#import-the-client-scope)Import the Client Scope

Use the following `curl` command to import the client scope into Keycloak:

```
curl -X POST "https://localhost:8443/admin/realms/oid4vc-vci/client-scopes" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d @client-scopes.json
```

- Replace `$ACCESS_TOKEN` with a valid **Keycloak Admin API access token**.
- **Avoid using `-k` in production**; instead, configure a **trusted TLS certificate**.
- If updating an existing scope, use `PUT` instead of `POST`.

### [](#create-the-client)Create the Client

Set up a client to handle Verifiable Credential (VC) requests and assign the necessary scopes. The client does not differ from regular OpenID Connect clients — with one exception: it must have the appropriate **optional ClientScopes** assigned that define the Verifiable Credentials it is allowed to issue.

1. Create a JSON file (e.g., `oid4vc-rest-api-client.json`) with the following content:
   
   ```
   {
     "clientId": "oid4vc-rest-api",
     "enabled": true,
     "protocol": "openid-connect",
     "publicClient": false,
     "serviceAccountsEnabled": true,
     "clientAuthenticatorType": "client-secret",
     "redirectUris": ["http://localhost:8080/*"],
     "directAccessGrantsEnabled": true,
     "defaultClientScopes": ["profile"],
     "optionalClientScopes": ["vc-scope-mapping"],
     "attributes": {
       "client.secret.creation.time": "1719785014",
       "client.introspection.response.allow.jwt.claim.enabled": "false",
       "login_theme": "keycloak",
       "post.logout.redirect.uris": "http://localhost:8080"
     }
   }
   ```
   
   - **clientId**: Unique identifier for the client.
   - **optionalClientScopes**: Links the `vc-scope-mapping` scope for VC requests.
2. Import the client using the following `curl` command:
   
   ```
   curl -k -X POST "https://localhost:8443/admin/realms/oid4vc-vci/clients" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d @oid4vc-rest-api-client.json
   ```

### [](#verify-the-configuration)Verify the Configuration

Validate the setup by accessing the **issuer metadata endpoint**:

1. Open a browser or use a tool like `curl` to visit:
   
   ```
   https://localhost:8443/.well-known/openid-credential-issuer/realms/oid4vc-vci
   ```

A successful response returns a JSON object containing details such as: - **Supported claims** - **Credential formats** - **Issuer metadata**

### [](#conclusion)Conclusion

You have successfully configured **Keycloak as a Verifiable Credential Issuer** using the **OID4VCI protocol**. This setup leverages Keycloak’s robust **identity management capabilities** to issue secure, **standards-compliant VCs**.

For a **complete reference implementation**, see the sample project: [Keycloak SSI Deployment](https://github.com/adorsys/Keycloak-ssi-deployment/tree/main).

## [](#_vault-administration)Using a vault to obtain secrets

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/vault.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fvault.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fvault.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak currently provides two out-of-the-box implementations of the Vault SPI: a plain-text file-based vault and Java KeyStore-based vault.

To obtain a secret from a vault rather than entering it directly, enter the following specially crafted string into the appropriate field:

```
${vault.key}
```

where the `key` is the name of the secret recognized by the vault.

To prevent secrets from leaking across realms, Keycloak combines the realm name with the `key` obtained from the vault expression. This method means that the `key` does not directly map to an entry in the vault but creates the final entry name according to the algorithm used to combine the `key` with the realm name. In case of the file-based vault, such combination reflects to a specific filename, for the Java KeyStore-based vault it’s a specific alias name.

You can obtain the secret from the vault in the following fields:

SMTP password

In the realm [SMTP settings](#_email)

LDAP bind credential

In the [LDAP settings](#_ldap) of LDAP-based user federation.

OIDC identity provider secret

In the *Client Secret* inside identity provider [OpenID Connect Config](#_identity_broker_oidc)

OIDC client secret

In the *Client Secret* inside confidential [OpenID Connect Client](#_client-credentials).

### [](#_vault-key-resolvers)Key resolvers

All built-in providers support the configuration of key resolvers. A key resolver implements the algorithm or strategy for combining the realm name with the key, obtained from the `${vault.key}` expression, into the final entry name used to retrieve the secret from the vault. Keycloak uses the `keyResolvers` property to configure the resolvers that the provider uses. The value is a comma-separated list of resolver names. An example of the configuration for the `files-plaintext` provider follows:

```
kc.[sh|bat] start --spi-vault--file--key-resolvers=REALM_UNDERSCORE_KEY,KEY_ONLY
```

The resolvers run in the same order you declare them in the configuration. For each resolver, Keycloak uses the last entry name the resolver produces, which combines the realm with the vault key to search for the vault’s secret. If Keycloak finds a secret, it returns the secret. If not, Keycloak uses the next resolver. This search continues until Keycloak finds a non-empty secret or runs out of resolvers. If Keycloak finds no secret, Keycloak returns an empty secret.

In the previous example, Keycloak uses the `REALM_UNDERSCORE_KEY` resolver first. If Keycloak finds an entry in the vault that using that resolver, Keycloak returns that entry. If not, Keycloak searches again using the `KEY_ONLY` resolver. If Keycloak finds an entry by using the `KEY_ONLY` resolver, Keycloak returns that entry. If Keycloak uses all resolvers, Keycloak returns an empty secret.

A list of the currently available resolvers follows:

  Name Description

KEY\_ONLY

Keycloak ignores the realm name and uses the key from the vault expression. Keycloak escapes occurrences of underscores in the key with another underscore character. For example, if the key is called `my_secret`, Keycloak searches for an entry in the vault named `my__secret`. This is to prevent conflicts with the default `REALM_UNDERSCORE_KEY` resolver.

REALM\_UNDERSCORE\_KEY

Keycloak combines the realm and key by using an underscore character. Keycloak escapes occurrences of underscores in the realm or key with another underscore character. For example, if the realm is called `master_realm` and the key is `smtp_key`, the combined key is `master__realm_smtp__key`.

REALM\_FILESEPARATOR\_KEY

Keycloak combines the realm and key by using the platform file separator character. The vault expression prohibits the use of characters that could cause path traversal, thus preventing access to secrets outside the corresponding realm.

FACTORY\_PROVIDED

Keycloak combines the realm and key by using the vault provider factory’s `VaultKeyResolver`, allowing the creation of a custom key resolver by extending an existing factory and implementing the `getFactoryResolver` method.

If you have not configured a resolver for the built-in providers, Keycloak selects the `REALM_UNDERSCORE_KEY`.

## [](#configuring-auditing-to-track-events)Configuring auditing to track events

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/events.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fevents.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fevents.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak includes a suite of auditing capabilities. You can record every login and administrator action and review those actions in the Admin Console. Keycloak also includes a Listener SPI that listens for events and can trigger actions. Examples of built-in listeners include log files and sending emails if an event occurs.

### [](#auditing-user-events)Auditing user events

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/events/login.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fevents%2Flogin.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fevents%2Flogin.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can record and view every event that affects users. Keycloak triggers login events for actions such as successful user login, a user entering an incorrect password, or a user account updating. By default, Keycloak does not store or display events in the Admin Console. Only the error events are logged to the Admin Console and the server’s log file.

Procedure

Use this procedure to start auditing user events.

1. Click **Realm settings** in the menu.
2. Click the **Events** tab.
3. Click the **User events settings** tab.
4. Toggle **Save events** to **ON**.
   
   User events settings
   
   ![User events settings](./images/user-events-settings.png)
5. Specify the length of time to store events in the **Expiration** field.
6. Click **Add saved types** to see other events you can save.
   
   Add types
   
   ![Add types](./images/add-event-types.png)
7. Click **Add**.

Click **Clear user events** when you want to delete all saved events.

Procedure

You can now view events.

1. Click the **Events** tab in the menu.
   
   User events
   
   ![Login Events](./images/user-events.png)
2. To filter events, click **Search user event**.
   
   Search user event
   
   ![Search user event](./images/search-user-event.png)

#### [](#event-types)Event types

**Login events:**

  Event Description

Login

A user logs in.

Register

A user registers.

Logout

A user logs out.

Code to Token

An application, or client, exchanges a code for a token.

Refresh Token

An application, or client, refreshes a token.

**Brute force protection:**

  Event Description

User disabled by permanent lockout

Brute force protection disabled the user account permanently due to too many login failures.

User disabled by temporary lockout

Brute force protection disabled the user account temporarily due to too many login failures.

**Identity Brokering:**

  Event Description

Federated identity link override

An existing Federated identity link was overridden

Federated identity link override error

Error occurred when trying to override an existing Federated identity link

**OAuth:**

  Event Description

OAuth2 extension grant

OAuth2 grant was executed

OAuth2 extension grant error

Error occurred during OAuth2 grant execution

**Account events:**

  Event Description

Social Link

A user account links to a social media provider.

Remove Social Link

The link from a social media account to a user account severs.

Update Email

An email address for an account changes.

Update Profile

A profile for an account changes.

Send Password Reset

Keycloak sends a password reset email.

Update Password (deprecated)

The password for an account changes.

Update Credential

The password or (time-based) one-time Password (OTP/TOTP) settings for an account changes.

Update TOTP (deprecated)

The Time-based One-time Password (TOTP) settings for an account changes.

Remove TOTP (deprecated)

Keycloak removes TOTP from an account.

Remove Credential

Keycloak removes a credential from an account.

Send Verify Email

Keycloak sends an email verification email.

Verify Email

Keycloak verifies the email address for an account.

Each event has a corresponding error event.

#### [](#event-listener)Event listener

Event listeners listen for events and perform actions based on that event. Keycloak includes two built-in listeners, the Logging Event Listener and Email Event Listener.

##### [](#the-logging-event-listener)The logging event listener

When the Logging Event Listener is enabled, this listener writes to a log file when an error event occurs.

An example log message from a Logging Event Listener:

```
11:36:09,965 WARN  [org.keycloak.events] (default task-51) type=LOGIN_ERROR, realmId=master,
                    clientId=myapp,
                    userId=19aeb848-96fc-44f6-b0a3-59a17570d374, ipAddress=127.0.0.1,
                    error=invalid_user_credentials, auth_method=openid-connect, auth_type=code,
                    redirect_uri=http://localhost:8180/myapp,
                    code_id=b669da14-cdbb-41d0-b055-0810a0334607, username=admin
```

You can use the Logging Event Listener to protect against hacker bot attacks:

1. Parse the log file for the `LOGIN_ERROR` event.
2. Extract the IP Address of the failed login event.
3. Send the IP address to an intrusion prevention software framework tool.

The Logging Event Listener logs events to the `org.keycloak.events` log category. Keycloak does not include debug log events in server logs, by default.

To include debug log events in server logs:

1. Change the log level for the `org.keycloak.events` category
2. Change the log level used by the Logging Event listener.

To change the log level used by the Logging Event listener, add the following:

```
bin/kc.[sh|bat] start --spi-events-listener--jboss-logging--success-level=info --spi-events-listener--jboss-logging--error-level=error
```

The valid values for log levels are `debug`, `info`, `warn`, `error`, and `fatal`.

##### [](#the-email-event-listener)The Email Event Listener

The Email Event Listener sends a message to the user’s email address when an event occurs and supports the following events:

- Login Error.
- Update Password.
- Update Time-based One-time Password (TOTP).
- Remove One-time Password (OTP).
- Update Credential.
- Remove Credential.

Below are the optional events you can configure:

- User disabled by permanent lockout.
- User disabled by temporary lockout.

The following conditions need to be met for an email to be sent:

- User has an email address.
- User’s email address is marked as verified.

Prerequisites

- Realm’s email settings configured.

Procedure

To enable the Email Listener:

1. Click **Realm settings** in the menu.
2. Click the **Events** tab.
3. Click the **Event listeners** field.
4. Select `email`.
   
   Event listeners
   
   ![Event listeners](./images/event-listeners.png)

You can exclude events by using the `--spi-events-listener--email--exclude-events` argument. For example:

```
kc.[sh|bat] --spi-events-listener--email--exclude-events=UPDATE_CREDENTIAL,REMOVE_CREDENTIAL
```

To enable optional events, use the following command:

```
kc.[sh|bat] --spi-events-listener--email--include-events=USER_DISABLED_BY_TEMPORARY_LOCKOUT_ERROR,USER_DISABLED_BY_PERMANENT_LOCKOUT
```

### [](#auditing-admin-events)Auditing admin events

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/events/admin.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fevents%2Fadmin.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fevents%2Fadmin.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

You can record all actions that are performed by an administrator in the Admin Console. The Admin Console performs administrative actions by invoking the Keycloak REST interface and Keycloak audits these REST invocations. You can view the resulting events in the Admin Console.

Procedure

Use this procedure to start auditing admin actions.

1. Click **Realm settings** in the menu.
2. Click the **Events** tab.
3. Click the **Admin events settings** tab.
4. Toggle **Save events** to **ON**.
   
   Keycloak displays the **Include representation** switch.
5. Toggle **Include representation** to **ON**.
   
   The `Include Representation` switch includes JSON documents sent through the admin REST API so you can view the administrators actions.
   
   Admin events settings
   
   ![Admin events settings](./images/admin-events-settings.png)
6. Click **Save**.
7. To clear the database of stored actions, click **Clear admin events**.

Procedure

You can now view admin events.

1. Click **Events** in the menu.
2. Click the **Admin events** tab.
   
   Admin events
   
   ![Admin events](./images/admin-events.png)

When the `Include Representation` switch is ON, it can lead to storing a lot of information in the database. You can set a maximum length of the representation by using the `--spi-events-store--jpa--max-field-length` argument. This setting is useful if you want to adhere to the underlying storage limitation. For example:

```
kc.[sh|bat] --spi-events-store--jpa--max-field-length=2500
```

## [](#mitigating_security_threats)Mitigating security threats

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Security vulnerabilities exist in any authentication server. See the Internet Engineering Task Force’s (IETF) [OAuth 2.0 Threat Model](https://datatracker.ietf.org/doc/html/rfc6819) and the [OAuth 2.0 Security Best Current Practice](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics) for more information.

### [](#host)Host

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/host.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fhost.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fhost.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak uses the public hostname in several ways, such as within token issuer fields and URLs in password reset emails.

By default, the hostname derives from request headers. No validation exists to ensure a hostname is valid. If you are not using a load balancer, or proxy, with Keycloak to prevent invalid host headers, configure the acceptable hostnames.

The hostname’s Service Provider Interface (SPI) provides a way to configure the hostname for requests. You can use this built-in provider to set a fixed URL for frontend requests while allowing backend requests based on the request URI. If the built-in provider does not have the required capability, you can develop a customized provider.

### [](#admin-endpoints-and-admin-console)Admin endpoints and Admin Console

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/admin.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fadmin.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fadmin.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak exposes the administrative REST API and the web console on the same port as non-administrative usage. Do not expose administrative endpoints externally if external access is not necessary.

### [](#password-guess-brute-force-attacks)Brute force attacks

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/brute-force.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fbrute-force.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fbrute-force.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

A brute force attack attempts to guess a user’s password by trying to log in multiple times. Keycloak has brute force detection capabilities and can permanently or temporarily disable a user account if the number of login failures exceeds a specified threshold.

The protection is applied only to authentication mechanisms prone to the brute force attacks. Keycloak applies it only to password, OTP and recovery codes.

When a user is locked and attempts to log in, Keycloak displays the default `Invalid username or password` error message. This message is the same error message as the message displayed for an invalid username or invalid password to ensure the attacker is unaware the account is disabled.

Brute force detection is disabled by default. Enable this feature to protect against brute force attacks.

To enable this protection:

1. Click **Realm Settings** in the menu
2. Click the **Security Defenses** tab.
3. Click the **Brute Force Detection** tab.
4. Choose the **Brute Force Mode** which best fit to your requirements.
   
   Brute force detection
   
   ![brute force](./images/brute-force.png)

#### [](#lockout-permanently)Lockout permanently

Keycloak disables a user account (blocking log in attempts) until an administrator re-enables it.

Lockout permanently

![brute force permanently](./images/brute-force-permanently.png)

**Permanent Lockout Parameters**

   Name Description Default

Max Login Failures

The maximum number of login failures.

30 failures

Quick Login Check Milliseconds

The minimum time between login attempts.

1000 milliseconds

Minimum Quick Login Wait

The minimum time the user is disabled when login attempts are quicker than *Quick Login Check Milliseconds*.

1 minute

**Permanent Lockout Algorithm**

1. On successful login
   
   1. Reset `count`
   2. call [`secondary authentication check (success)`](#_secondary-authn-check)
2. On failed login
   
   1. Increment `count`
   2. If `count` is greater than or equals to `Max login failures`
      
      1. locks the user
   3. Else if the time between this failure and the last failure is less than *Quick Login Check Milliseconds*
      
      1. Locks the user for the time specified at *Minimum Quick Login Wait*
   4. call [`secondary authentication check (fail)`](#_secondary-authn-check)

Enabling an user account resets the `count`.

#### [](#lockout-temporarily)Lockout temporarily

Keycloak disables a user account for a specific period of time. The time period that the account is disabled increases as the attack continues.

Lockout temporarily

![brute force temporarily](./images/brute-force-temporarily.png)

**Temporary Lockout Parameters**

   Name Description Default

Max Login Failures

The maximum number of login failures.

30 failures

Strategy to increase wait time

Strategy to increase the time a user will be temporarily disabled when the user’s login attempts exceed *Max Login Failures*

Multiple

Wait Increment

The time added to the time a user is temporarily disabled when the user’s login attempts exceed *Max Login Failures*.

1 minute

Max Wait

The maximum time a user is temporarily disabled.

15 minutes

Failure Reset Time

The time when the failure count resets. The timer runs from the last failed login. Make sure this number is always greater than `Max wait`; otherwise the effective wait time will never reach the value you have set to `Max wait`.

12 hours

Quick Login Check Milliseconds

The minimum time between login attempts.

1000 milliseconds

Minimum Quick Login Wait

The minimum time the user is disabled when login attempts are quicker than *Quick Login Check Milliseconds*.

1 minute

**Temporary Lockout Algorithm**

1. On successful login
   
   1. Reset `count`
   2. call [`secondary authentication check (success)`](#_secondary-authn-check)
2. On failed login
   
   1. If the time between this failure and the last failure is greater than *Failure Reset Time*
      
      1. Reset `count`
   2. Increment `count`
   3. Calculate `wait` according the brute force strategy defined (see below Strategies to set Wait Time).
   4. If `wait` is less than or equals to 0 and the time between this failure and the last failure is less than *Quick Login Check Milliseconds*
      
      1. set `wait` to *Minimum Quick Login Wait*
   5. if `wait` is greater than 0
      
      1. Temporarily disable the user for the smallest of `wait` and *Max Wait* seconds
   6. call [`secondary authentication check (fail)`](#_secondary-authn-check)

`count` does not increment when a temporarily disabled account commits a login failure.

**Strategies to set Wait Time**

Keycloak provides two strategies to calculate wait time: By multiples or Linear. By multiples is the first strategy introduced by Keycloak, so that is the default one.

By multiples strategy, wait time is incremented when the number (or count) of failures are multiples of `Max Login Failure`. For instance, if you set `Max Login Failures` to `5` and a `Wait Increment` to `30` seconds, the effective time that an account is disabled after several failed authentication attempts will be:

`Number of Failures`

`Wait Increment`

`Max Login Failures`

`Effective Wait Time`

1

30

5

0

2

30

5

0

3

30

5

0

4

30

5

0

**5**

**30**

5

**30**

6

30

5

30

7

30

5

30

8

30

5

30

9

30

5

30

**10**

**30**

5

**60**

At the fifth failed attempt, the account is disabled for `30` seconds. After reaching the next multiple of `Max Login Failures`, in this case `10`, the time increases from `30` to `60` seconds.

The By multiple strategy uses the following formula to calculate wait time: *Wait Increment in Seconds* * (`count` / *Max Login Failures*). The division is an integer division rounded down to a whole number.

For linear strategy, wait time is incremented when the `count` (or number) of failures is greater than or equals to `Max Login Failure`. For instance, if you have set `Max Login Failures` to `5` and a `Wait Increment` to\`30\` seconds, the effective time that an account is disabled after several failed authentication attempts will be:

`Number of Failures`

`Wait Increment`

`Max Login Failures`

`Effective Wait Time`

1

30

5

0

2

30

5

0

3

30

5

0

4

30

5

0

**5**

**30**

5

**30**

**6**

**30**

5

**60**

**7**

**30**

5

**90**

**8**

**30**

5

**120**

**9**

**30**

5

**150**

**10**

**30**

5

**180**

At the fifth failed attempt, the account is disabled for `30` seconds. Each new failure increases wait time according value specified at `wait increment`.

The linear strategy uses the following formula to calculate wait time: *Wait Increment in Seconds* * (1 + `count` - *Max Login Failures*).

#### [](#lockout-permanently-after-temporary-lockout)Lockout permanently after temporary lockout

Mixed mode. Locks user temporarily for specified number of times and then locks user permanently.

Lockout permanently after temporary lockout

![brute force mixed](./images/brute-force-mixed.png)

**Permanent lockout after temporary lockouts Parameters**

   Name Description Default

Max Login Failures

The maximum number of login failures.

30 failures

Maximum temporary Lockouts

The maximum number of temporary lockouts permitted before permanent lockout occurs.

1

Strategy to increase wait time

Strategy to increase the time a user will be temporarily disabled when the user’s login attempts exceed *Max Login Failures*

Multiple

Wait Increment

The time added to the time a user is temporarily disabled when the user’s login attempts exceed *Max Login Failures*.

1 minute

Max Wait

The maximum time a user is temporarily disabled.

15 minutes

Failure Reset Time

The time when the failure count resets. The timer runs from the last failed login. Make sure this number is always greater than `Max wait`; otherwise the effective wait time will never reach the value you have set to `Max wait`.

12 hours

Quick Login Check Milliseconds

The minimum time between login attempts.

1000 milliseconds

Minimum Quick Login Wait

The minimum time the user is disabled when login attempts are quicker than *Quick Login Check Milliseconds*.

1 minute

**Permanent lockout after temporary lockouts Algorithm**

1. On successful login
   
   1. Reset `count`
   2. Reset `temporary lockout` counter
   3. call [`secondary authentication check (success)`](#_secondary-authn-check)
2. On failed login
   
   1. If the time between this failure and the last failure is greater than *Failure Reset Time*
      
      1. Reset `count`
      2. Reset `temporary lockout` counter
   2. Increment `count`
   3. Calculate `wait` according the brute force strategy defined (see below Strategies to set Wait Time).
   4. If `wait` is less than or equals to 0 and the time between this failure and the last failure is less than *Quick Login Check Milliseconds*
      
      1. set `wait` to *Minimum Quick Login Wait*
      2. set `quick login failure` to ``true` ``
   5. if `wait` and `Maximum temporary Lockouts` is greater than 0
      
      1. set `wait` to the smallest of `wait` and *Max Wait* seconds
   6. if `quick login failure` is `false`
      
      1. Increment `temporary lockout` counter
   7. call [`secondary authentication check (fail)`](#_secondary-authn-check)
   8. If `temporary lockout` counter exceeds `Maximum temporary lockouts`
      
      1. Permanently locks the user
   9. Else
      
      1. Temporarily blocks the user according `wait` value

`count` does not increment when a temporarily disabled account commits a login failure.

#### [](#secondary-authentication-failures-lockout)Secondary Authentication Failures Lockout

When this lockout mechanism is enabled it is valid for all brute force modes. The mechanism is enabled when non zero value is set for `Maximum Secondary Authentication Failures`. It is disabled by default (`0`).

If the maximum number secondary authentication failures is exceeded the account is permanently locked. This behavior prevents attacks on second factor authenticators that can actually be brute forced (currently only OTP), when attacker already guessed the password. The reasonable value might be set for `100`.

**Secondary Authentication Check Algorithm**

1. If not OTP auth mechanism is used or secondary authentication check is disabled
   
   1. return
2. If called with `success` argument
   
   1. reset `Maximum Secondary Authentication Failures`
3. If called with `failure` argument
   
   1. Increment `Maximum Secondary Authentication Failures` counter
   2. If `Maximum Secondary Authentication Failures` exceeds desired number (100)
      
      1. Permanently locks the user

#### [](#downside-of-keycloak-brute-force-detection)Downside of Keycloak brute force detection

The downside of Keycloak brute force detection is that the server becomes vulnerable to denial of service attacks. When implementing a denial of service attack, an attacker can attempt to log in by guessing passwords for any accounts it knows and eventually causing Keycloak to disable the accounts.

Consider using intrusion prevention software (IPS). Keycloak logs every login failure and client IP address failure. You can point the IPS to the Keycloak server’s log file, and the IPS can modify firewalls to block connections from these IP addresses.

### [](#password-policies)Password policies

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/password.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fpassword.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fpassword.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Ensure you have a complex password policy to force users to choose complex passwords. See the [Password Policies](#_password-policies) chapter for more information. Prevent password guessing by setting up the Keycloak server to use one-time-passwords.

### [](#read_only_user_attributes)Read-only user attributes

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/read-only-attributes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fread-only-attributes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fread-only-attributes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Typical users who are stored in Keycloak have various attributes related to their user profiles. Such attributes include email, firstName or lastName. However users may also have attributes, which are not typical profile data, but rather metadata. The metadata attributes usually should be read-only for the users and the typical users never should have a way to update those attributes from the Keycloak user interface or Account REST API. Some of the attributes should be even read-only for the administrators when creating or updating user with the Admin REST API.

The metadata attributes are usually attributes from those groups:

- Various links or metadata related to the user storage providers. For example in case of the LDAP integration, the `LDAP_ID` attribute contains the ID of the user in the LDAP server.
- Metadata provisioned by User Storage. For example `createdTimestamp` provisioned from the LDAP should be always read-only by user or administrator.
- Metadata related to various authenticators. For example `KERBEROS_PRINCIPAL` attribute can contain the kerberos principal name of the particular user. Similarly attribute `usercertificate` can contain metadata related to binding the user with the data from the X.509 certificate, which is used typically when X.509 certificate authentication is enabled.
- Metadata related to the identificator of users by the applications/clients. For example `saml.persistent.name.id.for.my_app` can contain SAML NameID, which will be used by the client application `my_app` as the identifier of the user.
- Metadata related to the authorization policies, which are used for the attribute based access control (ABAC). Values of those attributes may be used for the authorization decisions. Hence it is important that those attributes cannot be updated by the users.

From the long term perspective, Keycloak will have a proper User Profile SPI, which will allow fine-grained configuration of every user attribute. Currently this capability is not fully available yet. So Keycloak has the internal list of user attributes, which are read-only for the users and read-only for the administrators configured at the server level.

This is the list of the read-only attributes, which are used internally by the Keycloak default providers and functionalities and hence are always read-only:

- For users: `KERBEROS_PRINCIPAL`, `LDAP_ID`, `LDAP_ENTRY_DN`, `CREATED_TIMESTAMP`, `createTimestamp`, `modifyTimestamp`, `userCertificate`, `saml.persistent.name.id.for.*`, `ENABLED`, `EMAIL_VERIFIED`
- For administrators: `KERBEROS_PRINCIPAL`, `LDAP_ID`, `LDAP_ENTRY_DN`, `CREATED_TIMESTAMP`, `createTimestamp`, `modifyTimestamp`

System administrators have a way to add additional attributes to this list. The configuration is currently available at the server level.

You can add this configuration by using the `spi-user-profile—​declarative-user-profile—​read-only-attributes` and `spi-user-profile—​declarative-user-profile—​admin-read-only-attributes` options. For example:

```
kc.[sh|bat] start --spi-user-profile--declarative-user-profile--read-only-attributes=foo,bar*
```

For this example, users and administrators would not be able to update attribute `foo`. Users would not be able to edit any attributes starting with the `bar`. So for example `bar` or `barrier`. Configuration is case-insensitive, so attributes like `FOO` or `BarRier` will be denied as well for this example. The wildcard character `*` is supported only at the end of the attribute name, so the administrator can effectively deny all the attributes starting with the specified character. The `*` in the middle of the attribute is considered as a normal character.

### [](#validate_user_attributes)Validate user attributes

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/validate-user-attributes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fvalidate-user-attributes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fvalidate-user-attributes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

With the functionality in [Managing user attributes](#user-profile), administrators can restrict the data users enter for attributes, for example, in user registration or the account console.

Administrators should not allow unmanaged attributes for users to prevent attackers adding an unlimited number of attributes. Attributes should have a validation that restricts the amount of data entered by attackers.

When using regular expressions to validate user attributes, avoid regular expressions that use an excessive amount of memory or CPU. See [OWASP’s Regular expression Denial of Service](https://owasp.org/www-community/attacks/Regular_expression_Denial_of_Service_-_ReDoS) for details.

### [](#clickjacking)Clickjacking

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/clickjacking.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fclickjacking.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fclickjacking.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Clickjacking is a technique of tricking users into clicking on a user interface element different from what users perceive. A malicious site loads the target site in a transparent iFrame, overlaid on top of a set of dummy buttons placed directly under important buttons on the target site. When a user clicks a visible button, they are clicking a button on the hidden page. An attacker can steal a user’s authentication credentials and access their resources by using this method.

By default, every response by Keycloak sets some specific HTTP headers that can prevent this from happening. Specifically, it sets [X-Frame-Options](https://datatracker.ietf.org/doc/html/rfc7034) and [Content-Security-Policy](https://www.w3.org/TR/CSP/). You should take a look at the definition of both of these headers as there is a lot of fine-grain browser access you can control.

Procedure

In the Admin Console, you can specify the values of the X-Frame-Options and Content-Security-Policy headers.

1. Click the **Realm Settings** menu item.
2. Click the **Security Defenses** tab.
   
   Security Defenses
   
   ![Security Defenses](./images/security-headers.png)

By default, Keycloak only sets up a *same-origin* policy for iframes.

### [](#sslhttps-requirement)SSL/HTTPS requirement

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/ssl.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fssl.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fssl.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

OAuth 2.0/OpenID Connect uses access tokens for security. Attackers can scan your network for access tokens and use them to perform malicious operations for which the token has permission. This attack is known as a man-in-the-middle attack. Use SSL/HTTPS for communication between the Keycloak auth server and the clients Keycloak secures to prevent man-in-the-middle attacks.

Keycloak has [three modes for SSL/HTTPS](#_ssl_modes). SSL is complex to set up, so Keycloak allows non-HTTPS communication over private IP addresses such as localhost, 192.168.x.x, and other private IP addresses. In production, ensure you enable SSL and SSL is compulsory for all operations.

On the adapter/client-side, you can disable the SSL trust manager. The trust manager ensures the client’s identity that Keycloak communicates with is valid and ensures the DNS domain name against the server’s certificate. In production, ensure that each of your client adapters uses a truststore to prevent DNS man-in-the-middle attacks.

### [](#csrf-attacks)CSRF attacks

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/csrf.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fcsrf.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fcsrf.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

A Cross-site request forgery (CSRF) attack uses HTTP requests from users that websites have already authenticated. Any site using cookie-based authentication is vulnerable to CSRF attacks. You can mitigate these attacks by matching a state cookie against a posted form or query parameter.

The OAuth 2.0 login specification requires that a state cookie matches against a transmitted state parameter. Keycloak fully implements this part of the specification, so all logins are protected.

The Keycloak Admin Console is a JavaScript/HTML5 application that makes REST calls to the backend Keycloak admin REST API. These calls all require bearer token authentication and consist of JavaScript Ajax calls, so CSRF is impossible. You can configure the admin REST API to validate the CORS origins.

The Account Console in Keycloak can be vulnerable to CSRF. To prevent CSRF attacks, Keycloak sets a state cookie and embeds the value of this cookie in hidden form fields or query parameters within action links. Keycloak checks the query/form parameter against the state cookie to verify that the same user made the call.

### [](#unspecific-redirect-uris_server_administration_guide)Unspecific redirect URIs

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/redirect.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fredirect.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fredirect.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Make your registered redirect URIs as specific as feasible. Registering vague redirect URIs for [Authorization Code Flows](#con-oidc-auth-flows_server_administration_guide) can allow malicious clients to impersonate another client with broader access. Impersonation can happen if two clients live under the same domain, for example.

You can use secure redirect uris enforcer executor for your realm. The result makes sure that client administrators are able to register only clients with specific redirect-uris matching various requirements such as requiring that a URL cannot have wildcards in the context path or can be limited to specified permitted domains. See [Client Policies](#_client_policies) for details about how to configure client policies with a specific executor.

### [](#fapi-compliance)FAPI compliance

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/fapi-compliance.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Ffapi-compliance.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Ffapi-compliance.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To make sure that Keycloak server will validate your client to be more secure and FAPI compliant, you can configure client policies for the FAPI support. **FAPI** details are described in the [securing apps](https://www.keycloak.org/guides#securing-apps) section. Among other things, this ensures some security best practices described above like SSL required for clients, secure redirect URI used and more of similar best practices.

### [](#oauth-2-1-compliance)OAuth 2.1 compliance

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/oauth21-compliance.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Foauth21-compliance.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Foauth21-compliance.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

To make sure that Keycloak server will validate your client to be more secure and OAuth 2.1 compliant, you can configure client policies for the OAuth 2.1 support. **OAuth 2.1** details are described in the [securing apps](https://www.keycloak.org/guides#securing-apps) section.

### [](#compromised-access-and-refresh-tokens)Compromised access and refresh tokens

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/compromised-tokens.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fcompromised-tokens.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fcompromised-tokens.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak includes several actions to prevent malicious actors from stealing access tokens and refresh tokens. The crucial action is to enforce SSL/HTTPS communication between Keycloak and its clients and applications. Keycloak does not enable SSL by default.

Another action to mitigate damage from leaked access tokens is to shorten the token’s lifespans. You can specify token lifespans within the [timeouts page](#_timeouts). Short lifespans for access tokens force clients and applications to refresh their access tokens after a short time. If an admin detects a leak, the admin can log out all user sessions to invalidate these refresh tokens or set up a revocation policy.

Ensure refresh tokens always stay private to the client and are never transmitted.

You can mitigate damage from leaked access tokens and refresh tokens by issuing these tokens as holder-of-key tokens. See [OAuth 2.0 Mutual TLS Client Certificate Bound Access Token](#_mtls-client-certificate-bound-tokens) for more information.

If an access token or refresh token is compromised, access the Admin Console and push a not-before revocation policy to all applications. Pushing a not-before policy ensures that any tokens issued before that time become invalid. Pushing a new not-before policy ensures that applications must download new public keys from Keycloak and mitigate damage from a compromised realm signing key. See the [keys chapter](#realm_keys) for more information.

You can disable specific applications, clients, or users if they are compromised.

### [](#compromised-authorization-code)Compromised authorization code

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/compromised-codes.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fcompromised-codes.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fcompromised-codes.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

For the [OIDC Auth Code Flow](#con-oidc-auth-flows_server_administration_guide), Keycloak generates a cryptographically strong random value for its authorization codes. An authorization code is used only once to obtain an access token.

On the [timeouts page](#_timeouts) in the Admin Console, you can specify the length of time an authorization code is valid. Ensure that the length of time is less than 10 seconds, which is long enough for a client to request a token from the code.

You can also defend against leaked authorization codes by applying [Proof Key for Code Exchange (PKCE)](#_proof-key-for-code-exchange) to clients.

### [](#open-redirectors)Open redirectors

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/open-redirect.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fopen-redirect.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fopen-redirect.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

An open redirector is an endpoint using a parameter to automatically redirect a user agent to the location specified by the parameter value without validation. An attacker can use the end-user authorization endpoint and the redirect URI parameter to use the authorization server as an open redirector, using a user’s trust in an authorization server to launch a phishing attack.

Keycloak requires that all registered applications and clients register at least one redirection URI pattern. When a client requests that Keycloak performs a redirect, Keycloak checks the redirect URI against the list of valid registered URI patterns. Clients and applications must register as specific a URI pattern as possible to mitigate open redirector attacks.

If an application requires a non http(s) custom scheme, it should be an explicit part of the validation pattern (for example `custom:/app/*`). For security reasons a general pattern like `*` does not cover non http(s) schemes.

By using [Client Policies](#_client_policies), an administrator can make sure that clients cannot register open redirect URLs such as `*`.

### [](#ssrf)Mitigating Server-Side Request Forgery (SSRF)

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/ssrf.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fssrf.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fssrf.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

**Server-Side Request Forgery (SSRF)** is a security vulnerability where an attacker induces the server (Keycloak, in this case) to make requests to an unintended destination. In Keycloak, this risk is primarily associated with client fields that trigger requests from the server, such as the **JWKS URI**.

If these fields are not strictly validated, a malicious administrator, a compromised account or a Dynamic Client Registration could point them toward internal network resources or loopback interfaces (e.g., `http://localhost:8080/any`), allowing the attacker to probe or attack the internal infrastructure from the Keycloak server itself.

To mitigate this threat, Keycloak provides a solution using [Client Policies](#_client_policies) with some executors like **Secure Client Uris** to provide baseline protection for transport security and wildcards, but especially with the **Secure Client URIs Pattern** executor.

#### [](#secure-client-uris-pattern-executor)Secure Client URIs Pattern Executor

This executor enforces a strict security policy by validating client URIs against an **allowlist of patterns**. If a URI does not match at least one of the configured patterns, the client creation or update is rejected.

It can be used not only for `JWKS URI` but to validate all available uri fields of the client like including `rootUrl`, `adminUrl`, `redirectUris`, `webOrigins` and others.

Table 9. Configuration Properties   Configuration Description

Allowed URI Patterns

A list of Regular Expressions. A client URI is considered valid **only** if it matches at least one of these patterns. If this list is empty or invalid, the executor blocks **all** URIs.

Client URI Fields to validate

A list of specific client fields to validate (e.g., `jwksUri`, `adminUrl`). If left empty, **all** supported URI fields are validated by default. The list is available in the configuration.

#### [](#anti-ssrf-pattern-examples)Anti-SSRF Pattern Examples

To effectively prevent SSRF, Regular Expressions should be designed to exclude IP addresses and local hostnames.

**1. Restricting to Trusted Domains**

This pattern ensures that the Keycloak server only communicates with official domains, preventing the use of `localhost` or numerical IP addresses.

- **Pattern:** `^https://([a-zA-Z0-9-]+\.)?mycompany\.com(/.*)?$`
- **Effect:** Blocks attempts to use loopback addresses (e.g., `http://localhost:8080`) or cloud metadata endpoints.

**2. Enforcing HTTPS**

SSRF often exploits unencrypted protocols to bypass internal security controls. Enforcing HTTPS via patterns reduces the attack surface.

- **Pattern:** `^https://[a-zA-Z0-9.-]+(/.*)?$`
- **Result:** Rejects the use of the `http://` scheme, which is frequently used to probe internal services.

https can be enforced also with the **Secure Client Uris** executor but only together with wildcards and for all the client URIs

**3. Specific Provider Allowlisting**

If your clients use JWKS URI from third-party providers, only add the specific domains belonging to those providers.

- **Pattern:** `^https://trusted-provider\.io/.*`
- **Result:** All other domains will be rejected for JWKS requests.

Regular Expression Syntax

This executor uses Java Regular Expressions. Ensure you correctly escape special characters (for example, use `\.` to indicate a literal dot). An incorrect pattern could prevent legitimate updates to clients. It is recommended to test regexes in a non-production environment.

### [](#password-database-compromised)Password database compromised

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/password-db-compromised.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fpassword-db-compromised.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fpassword-db-compromised.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak does not store passwords in raw text but as hashed text, using the `PBKDF2-HMAC-SHA512` message digest algorithm. Keycloak performs `210,000` hashing iterations, the number of iterations recommended by the security community. This number of hashing iterations can adversely affect performance as PBKDF2 hashing uses a significant amount of CPU resources.

### [](#limiting-scope)Limiting scope

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/scope.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fscope.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fscope.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

#### [](#scope-availability)Scope availability

By default, new client applications have unlimited `role scope mappings`. Every access token for that client contains all permissions that the user has. If an attacker compromises the client and obtains the client’s access tokens, each system that the user can access is compromised.

Limit the roles of an access token by using the [Scope menu](#_role_scope_mappings) for each client. Alternatively, you can set role scope mappings at the Client Scope level and assign Client Scopes to your client by using the [Client Scope menu](#_client_scopes_linking).

Removing the offline scope for a client also removes the ability to issue long-lived offline tokens for a client and offers better control over sessions by users.

#### [](#scope-visibility)Scope visibility

By default, all scopes are included in the OpenID Connect discovery endpoint. To reduce the discoverability and OSINT-exposure, you can configure each scope to be excluded by disabling **Include in OpenID Provider Metadata**.

### [](#limit-token-audience)Limit token audience

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/audience-limit.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Faudience-limit.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Faudience-limit.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

In environments with low levels of trust among services, limit the audiences on the token. See the [OAuth2 Threat Model](https://datatracker.ietf.org/doc/html/rfc6819#section-5.1.5.5) and the [Audience Support](#audience-support) section for more information.

### [](#_limit-authentication-sessions)Limit Authentication Sessions

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/auth-sessions-limit.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fauth-sessions-limit.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fauth-sessions-limit.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

[Authentication sessions](#_authentication-sessions) track the state of the authentication. The text below is applicable regardless of the source flow.

This section describes deployments that use the Infinispan provider for authentication sessions.

Authentication session is internally stored as `RootAuthenticationSessionEntity`. Each `RootAuthenticationSessionEntity` can have multiple authentication sub-sessions stored within the `RootAuthenticationSessionEntity` as a collection of `AuthenticationSessionEntity` objects. Keycloak stores authentication sessions in a dedicated Infinispan cache. The number of `AuthenticationSessionEntity` per `RootAuthenticationSessionEntity` contributes to the size of each cache entry. Total memory footprint of authentication session cache is determined by the number of stored `RootAuthenticationSessionEntity` and by the number of `AuthenticationSessionEntity` within each `RootAuthenticationSessionEntity`.

The number of maintained `RootAuthenticationSessionEntity` objects corresponds to the number of unfinished login flows from the browser. To keep the number of `RootAuthenticationSessionEntity` under control, using an advanced firewall control to limit ingress network traffic is recommended.

Higher memory usage may occur for deployments where there are many active `RootAuthenticationSessionEntity` with a lot of `AuthenticationSessionEntity`. If the load balancer does not support or is not configured for session stickiness, the load over network in a cluster can increase significantly. The reason for this load is that each request that lands on a node that does not own the appropriate authentication session needs to retrieve and update the authentication session record in the owner node which involves a separate network transmission for both the retrieval and the storage.

The maximum number of `AuthenticationSessionEntity` per `RootAuthenticationSessionEntity` can be configured in `authenticationSessions` SPI by setting property `authSessionsLimit`. The default value is set to 300 `AuthenticationSessionEntity` per a `RootAuthenticationSessionEntity`. When this limit is reached, the oldest authentication sub-session will be removed after a new authentication session request.

The following example shows how to limit the number of active `AuthenticationSessionEntity` per a `RootAuthenticationSessionEntity` to 100.

```
bin/kc.[sh|bat] start --spi-authentication-sessions--infinispan--auth-sessions-limit=100
```

The equivalent command for the new map storage:

```
bin/kc.[sh|bat] start --spi-authentication-sessions--map--auth-sessions-limit=100
```

### [](#sql-injection-attacks)SQL injection attacks

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/threat/sql.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fthreat%2Fsql.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fthreat%2Fsql.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Currently, Keycloak has no known SQL injection vulnerabilities.

## [](#_account-service)Account Console

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/account.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Faccount.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Faccount.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

Keycloak users can manage their accounts through the Account Console. They can configure their profiles, add two-factor authentication, include identity provider accounts, and oversee device activity.

Additional resources

- The Account Console can be configured in terms of appearance and language preferences. An example is adding additional attributes to the **Personal info** page. For more information, see the [Server Developer Guide](https://www.keycloak.org/docs/26.6.3/server_development/).

### [](#accessing-the-account-console)Accessing the Account Console

Procedure

1. Make note of the realm name and IP address for the Keycloak server where your account exists.
2. In a web browser, enter a URL in this format: *server-root*/realms/{realm-name}/account.
3. Enter your login name and password.

Account Console

![Account Console](./images/account-console-intro.png)

You can also ask for additional scopes when calling the account console URL by setting the `scope` parameter in this format: *server-root*/realms/{realm-name}/account?scope=phone.

### [](#configuring-ways-to-sign-in)Configuring ways to sign in

You can sign in to this console using basic authentication (a login name and password) or two-factor authentication. For two-factor authentication, use one of the following procedures.

#### [](#two-factor-authentication-with-otp)Two-factor authentication with OTP

Prerequisites

- OTP is a valid authentication mechanism for your realm.

Procedure

1. Click **Account security** in the menu.
2. Click **Signing in**.
3. Click **Set up Authenticator application**.
   
   Signing in
   
   ![Signing in](./images/account-console-signing-in.png)
4. Follow the directions that appear on the screen to use your mobile device as your OTP generator.
5. Scan the QR code in the screen shot into the OTP generator on your mobile device.
6. Log out and log in again.
7. Respond to the prompt by entering an OTP that is provided on your mobile device.

#### [](#two-factor-authentication-with-webauthn)Two-factor authentication with WebAuthn

Prerequisites

- WebAuthn is a valid two-factor authentication mechanism for your realm. Please follow the [WebAuthn](#webauthn_server_administration_guide) section for more details.

Procedure

1. Click **Account Security** in the menu.
2. Click **Signing In**.
3. Click **Set up a Passkey**.
   
   Signing In
   
   ![Signing in with a Passkey](./images/account-console-signing-in-webauthn-2factor.png)
4. Prepare your Passkey. How you prepare this key depends on the type of Passkey you use. For example, for a USB based Yubikey, you may need to put your key into the USB port on your laptop.
5. Click **Register** to register your Passkey.
6. Log out and log in again.
7. Assuming authentication flow was correctly set, a message appears asking you to authenticate with your Passkey as second factor.

#### [](#passwordless-authentication-with-webauthn)Passwordless authentication with WebAuthn

Prerequisites

- WebAuthn is a valid passwordless authentication mechanism for your realm. Please follow the [Passwordless WebAuthn section](#_webauthn_passwordless) for more details.

Procedure

1. Click **Account Security** in the menu.
2. Click **Signing In**.
3. Click **Set up a Passkey** in the **Passwordless** section.
   
   Signing In
   
   ![Signing in with a Passkey](./images/account-console-signing-in-webauthn-passwordless.png)
4. Prepare your Passkey. How you prepare this key depends on the type of Passkey you use. For example, for a USB based Yubikey, you may need to put your key into the USB port on your laptop.
5. Click **Register** to register your Passkey.
6. Log out and log in again.
7. Assuming authentication flow was correctly set, a message appears asking you to authenticate with your Passkey as second factor. You no longer need to provide your password to log in.

### [](#viewing-device-activity)Viewing device activity

You can view the devices that are logged in to your account.

Procedure

1. Click **Account security** in the menu.
2. Click **Device activity**.
3. Log out a device if it looks suspicious.

Devices

![Devices](./images/account-console-device.png)

### [](#adding-an-identity-provider-account)Adding an identity provider account

You can link your account with an [identity broker](#_identity_broker). This option is often used to link social provider accounts.

Procedure

1. Log into the Admin Console.
2. Click **Identity providers** in the menu.
3. Select a provider and complete the fields.
4. Return to the Account Console.
5. Click **Account security** in the menu.
6. Click **Linked accounts**.

The identity provider you added appears in this page.

Linked Accounts

![Linked Accounts](./images/account-console-linked.png)

### [](#accessing-other-applications)Accessing other applications

The **Applications** menu item shows users which applications you can access. In this case, only the Account Console is available.

Applications

![Applications](./images/account-console-applications.png)

### [](#viewing-group-memberships)Viewing group memberships

You can view the groups you are associated with by clicking the **Groups** menu. If you select **Direct membership** checkbox, you will see only the groups you are direct associated with.

Prerequisites

- You need to have the **view-groups** account role for being able to view **Groups** menu.

View group memberships

![View group memberships](./images/account-console-groups.png)

## [](#admin-cli)Admin CLI

[Edit this section](https://github.com/keycloak/keycloak/tree/main/docs/documentation/server_admin/topics/admin-cli.adoc) [Report an issue](https://github.com/keycloak/keycloak/issues/new?template=bug.yml&title=Docs%3A%20server_admin%2Ftopics%2Fadmin-cli.adoc&description=%0A%0AFile%3A%20server_admin%2Ftopics%2Fadmin-cli.adoc&version=26.6.3&behaviorExpected=%3C%21--%20describe%20what%20you%20want%20to%20see%20in%20the%20docs%20--%3E&behaviorActual=%3C%21--%20describe%20what%20is%20currently%20wrong%20or%20missing%20in%20the%20docs%20--%3E&reproducer=%3C%21--%20list%20steps%20in%20the%20application%20that%20show%20behavior%20that%20should%20be%20documented%20--%3E)

With Keycloak, you can perform administration tasks from the command-line interface (CLI) by using the Admin CLI command-line tool.

### [](#installing-the-admin-cli)Installing the Admin CLI

Keycloak packages the Admin CLI server distribution with the execution scripts in the `bin` directory.

The Linux script is called `kcadm.sh`, and the script for Windows is called `kcadm.bat`. Add the Keycloak server directory to your `PATH` to use the client from any location on your file system.

For example:

- Linux:
  
  ```
  $ export PATH=$PATH:$KEYCLOAK_HOME/bin
  $ kcadm.sh
  ```
- Windows:
  
  ```
  c:\> set PATH=%PATH%;%KEYCLOAK_HOME%\bin
  c:\> kcadm
  ```

You must set the `KEYCLOAK_HOME` environment variable to the path where you extracted the Keycloak Server distribution.

To avoid repetition, the rest of this document only uses Windows examples in places where the CLI differences are more than just in the `kcadm` command name.

### [](#using-the-admin-cli)Using the Admin CLI

The Admin CLI makes HTTP requests to Admin REST endpoints. Access to the Admin REST endpoints requires authentication.

Consult the Admin REST API documentation for details about JSON attributes for specific endpoints.

1. Start an authenticated session by logging in. You can now perform create, read, update, and delete (CRUD) operations.
   
   For example:
   
   - Linux:
     
     ```
     $ kcadm.sh config credentials --server http://localhost:8080 --realm demo --user admin --client admin
     $ kcadm.sh create realms -s realm=demorealm -s enabled=true -o
     $ CID=$(kcadm.sh create clients -r demorealm -s clientId=my_client -s 'redirectUris=["http://localhost:8980/myapp/*"]' -i)
     $ kcadm.sh get clients/$CID/installation/providers/keycloak-oidc-keycloak-json
     ```
   - Windows:
     
     ```
     c:\> kcadm config credentials --server http://localhost:8080 --realm demo --user admin --client admin
     c:\> kcadm create realms -s realm=demorealm -s enabled=true -o
     c:\> kcadm create clients -r demorealm -s clientId=my_client -s "redirectUris=[\"http://localhost:8980/myapp/*\"]" -i > clientid.txt
     c:\> set /p CID=<clientid.txt
     c:\> kcadm get clients/%CID%/installation/providers/keycloak-oidc-keycloak-json
     ```
2. In a production environment, access Keycloak by using `https:` to avoid exposing tokens. If a trusted certificate authority, included in Java’s default certificate truststore, has not issued a server’s certificate, prepare a `truststore.jks` file and instruct the Admin CLI to use it.
   
   For example:
   
   - Linux:
     
     ```
     $ kcadm.sh config truststore --trustpass $PASSWORD ~/.keycloak/truststore.jks
     ```
   - Windows:
     
     ```
     c:\> kcadm config truststore --trustpass %PASSWORD% %HOMEPATH%\.keycloak\truststore.jks
     ```

### [](#sensitive-options)Sensitive Options

Sensitive values, such as passwords, may be specified as command options. That is generally not recommended. There are also mechanisms by which you can be prompted for the sensitive value by either omitting the option or providing a value. Finally all will have a corresponding env variable that can be used instead. Check the help of the command you are running to see all possible options.

### [](#authenticating)Authenticating

When you log in with the Admin CLI, you specify:

- A server endpoint URL
- A realm
- A user name

Another option is to specify a clientId only, which creates a unique service account for you to use.

When you log in using a user name, use a password for the specified user. When you log in using a clientId, you need the client secret only, not the user password. You can also use the `Signed JWT` rather than the client secret.

Ensure the account used for the session has the proper permissions to invoke Admin REST API operations. For example, the `realm-admin` role of the `realm-management` client can administer the realm of the user.

Two primary mechanisms are available for authentication. One mechanism uses `kcadm config credentials` to start an authenticated session.

```
$ kcadm.sh config credentials --server http://localhost:8080 --realm master --user admin
```

This mechanism maintains an authenticated session between the `kcadm` command invocations by saving the obtained access token and its associated refresh token. It can maintain other secrets in a private configuration file. See the [next chapter](#_working_with_alternative_configurations) for more information.

The second mechanism authenticates each command invocation for the duration of the invocation. This mechanism increases the load on the server and the time spent on round trips obtaining tokens. The benefit of this approach is that it is unnecessary to save tokens between invocations, so nothing is saved to disk. Keycloak uses this mode when the `--no-config` argument is specified.

For example, when performing an operation, specify all the information required for authentication.

```
$ kcadm.sh get realms --no-config --server http://localhost:8080 --realm master --user admin
```

Run the `kcadm.sh help` command for more information on using the Admin CLI.

Run the `kcadm.sh config credentials --help` command for more information about starting an authenticated session.

If you do not specify the --password option (it is generally recommended to not provide passwords as part of the command), you will be prompted for a password unless one is specified as the environment variable KC\_CLI\_PASSWORD.

### [](#_working_with_alternative_configurations)Working with alternative configurations

By default, the Admin CLI maintains a configuration file named `kcadm.config`. Keycloak places this file in the user’s home directory. In Linux-based systems, the full pathname is `$HOME/.keycloak/kcadm.config`. In Windows, the full pathname is `%HOMEPATH%\.keycloak\kcadm.config`.

You can use the `--config` option to point to a different file or location so you can maintain multiple authenticated sessions in parallel.

Perform operations tied to a single configuration file from a single thread.

Ensure the configuration file is invisible to other users on the system. It contains access tokens and secrets that must be private. Keycloak creates the `~/.keycloak` directory and its contents automatically with proper access limits. If the directory already exists, Keycloak does not update the directory’s permissions.

It is possible to avoid storing secrets inside a configuration file, but doing so is inconvenient and increases the number of token requests. Use the `--no-config` option with all commands and specify the authentication information the `config credentials` command requires with each invocation of `kcadm`.

### [](#basic-operations-and-resource-uris)Basic operations and resource URIs

The Admin CLI can generically perform CRUD operations against Admin REST API endpoints with additional commands that simplify particular tasks.

The main usage pattern is listed here:

```
$ kcadm.sh create ENDPOINT [ARGUMENTS]
$ kcadm.sh get ENDPOINT [ARGUMENTS]
$ kcadm.sh update ENDPOINT [ARGUMENTS]
$ kcadm.sh delete ENDPOINT [ARGUMENTS]
```

The `create`, `get`, `update`, and `delete` commands map to the HTTP verbs `POST`, `GET`, `PUT`, and `DELETE`, respectively. ENDPOINT is a target resource URI and can be absolute (starting with `http:` or `https:`) or relative, that Keycloak uses to compose absolute URLs in the following format:

```
SERVER_URI/admin/realms/REALM/ENDPOINT
```

For example, if you authenticate against the server [http://localhost:8080](http://localhost:8080) and realm is `master`, using `users` as ENDPOINT creates the [http://localhost:8080/admin/realms/master/users](http://localhost:8080/admin/realms/master/users) resource URL.

If you set ENDPOINT to `clients`, the effective resource URI is [http://localhost:8080/admin/realms/master/clients](http://localhost:8080/admin/realms/master/clients).

Keycloak has a `realms` endpoint that is the container for realms. It resolves to:

```
SERVER_URI/admin/realms
```

Keycloak has a `serverinfo` endpoint. This endpoint is independent of realms.

When you authenticate as a user with realm-admin powers, you may need to perform commands on multiple realms. If so, specify the `-r` option to tell the CLI which realm the command is to execute against explicitly. Instead of using `REALM` as specified by the `--realm` option of `kcadm.sh config credentials`, the command uses `TARGET_REALM`.

```
SERVER_URI/admin/realms/TARGET_REALM/ENDPOINT
```

For example:

```
$ kcadm.sh config credentials --server http://localhost:8080 --realm master --user admin
$ kcadm.sh create users -s username=testuser -s enabled=true -r demorealm
```

In this example, you start a session authenticated as the `admin` user in the `master` realm. You then perform a POST call against the resource URL `http://localhost:8080/admin/realms/demorealm/users`.

The `create` and `update` commands send a JSON body to the server. You can use `-f FILENAME` to read a pre-made document from a file. When you can use the `-f -` option, Keycloak reads the message body from the standard input. You can specify individual attributes and their values, as seen in the `create users` example. Keycloak composes the attributes into a JSON body and sends them to the server.

The value in name=value pairs used in --set, -s options, are assumed to be JSON. If it cannot be parsed as valid JSON, then it will be sent to the server as a text value.

If the value is enclosed in quotes after shell processing, but is not valid JSON, the quotes will be stripped and the rest of the value will be sent as text. This behavior is deprecated, please consider specifying your value without quotes or a valid JSON string literal with double quotes.

Several methods are available in Keycloak to update a resource using the `update` command. You can determine the current state of a resource and save it to a file, edit that file, and send it to the server for an update.

For example:

```
$ kcadm.sh get realms/demorealm > demorealm.json
$ vi demorealm.json
$ kcadm.sh update realms/demorealm -f demorealm.json
```

This method updates the resource on the server with the attributes in the sent JSON document.

Another method is to perform an on-the-fly update by using the `-s, --set` options to set new values.

For example:

```
$ kcadm.sh update realms/demorealm -s enabled=false
```

This method sets the `enabled` attribute to `false`.

By default, the `update` command performs a `get` and then merges the new attribute values with existing values. In some cases, the endpoint may support the `put` command but not the `get` command. You can use the `-n` option to perform a no-merge update, which performs a `put` command without first running a `get` command.

### [](#realm-operations)Realm operations

#### Creating a new realm

Use the `create` command on the `realms` endpoint to create a new enabled realm. Set the attributes to `realm` and `enabled`.

```
$ kcadm.sh create realms -s realm=demorealm -s enabled=true
```

Keycloak disables realms by default. You can use a realm immediately for authentication by enabling it.

A description for a new object can also be in JSON format.

```
$ kcadm.sh create realms -f demorealm.json
```

You can send a JSON document with realm attributes directly from a file or pipe the document to standard input.

For example:

- Linux:
  
  ```
  $ kcadm.sh create realms -f - << EOF
  { "realm": "demorealm", "enabled": true }
  EOF
  ```
- Windows:
  
  ```
  c:\> echo { "realm": "demorealm", "enabled": true } | kcadm create realms -f -
  ```

#### Listing existing realms

This command returns a list of all realms.

```
$ kcadm.sh get realms
```

Keycloak filters the list of realms on the server to return realms a user can see only.

The list of all realm attributes can be verbose, and most users are interested in a subset of attributes, such as the realm name and the enabled status of the realm. You can specify the attributes to return by using the `--fields` option.

```
$ kcadm.sh get realms --fields realm,enabled
```

You can display the result as comma-separated values.

```
$ kcadm.sh get realms --fields realm --format csv --noquotes
```

#### Getting a specific realm

Append a realm name to a collection URI to get an individual realm.

```
$ kcadm.sh get realms/master
```

#### Updating a realm

1. Use the `-s` option to set new values for the attributes when you do not want to change all of the realm’s attributes.
   
   For example:
   
   ```
   $ kcadm.sh update realms/demorealm -s enabled=false
   ```
2. If you want to set all writable attributes to new values:
   
   1. Run a `get` command.
   2. Edit the current values in the JSON file.
   3. Resubmit.
      
      For example:
      
      ```
      $ kcadm.sh get realms/demorealm > demorealm.json
      $ vi demorealm.json
      $ kcadm.sh update realms/demorealm -f demorealm.json
      ```

#### Deleting a realm

Run the following command to delete a realm:

```
$ kcadm.sh delete realms/demorealm
```

#### Turning on all login page options for the realm

Set the attributes that control specific capabilities to `true`.

For example:

```
$ kcadm.sh update realms/demorealm -s registrationAllowed=true -s registrationEmailAsUsername=true -s rememberMe=true -s verifyEmail=true -s resetPasswordAllowed=true -s editUsernameAllowed=true
```

#### Listing the realm keys

Use the `get` operation on the `keys` endpoint of the target realm.

```
$ kcadm.sh get keys -r demorealm
```

#### Generating new realm keys

1. Get the ID of the target realm before adding a new RSA-generated key pair.
   
   For example:
   
   ```
   $ kcadm.sh get realms/demorealm --fields id --format csv --noquotes
   ```
2. Add a new key provider with a higher priority than the existing providers as revealed by `kcadm.sh get keys -r demorealm`.
   
   For example:
   
   - Linux:
     
     ```
     $ kcadm.sh create components -r demorealm -s name=rsa-generated -s providerId=rsa-generated -s providerType=org.keycloak.keys.KeyProvider -s parentId=959844c1-d149-41d7-8359-6aa527fca0b0 -s 'config.priority=["101"]' -s 'config.enabled=["true"]' -s 'config.active=["true"]' -s 'config.keySize=["2048"]'
     ```
   - Windows:
     
     ```
     c:\> kcadm create components -r demorealm -s name=rsa-generated -s providerId=rsa-generated -s providerType=org.keycloak.keys.KeyProvider -s parentId=959844c1-d149-41d7-8359-6aa527fca0b0 -s "config.priority=[\"101\"]" -s "config.enabled=[\"true\"]" -s "config.active=[\"true\"]" -s "config.keySize=[\"2048\"]"
     ```
3. Set the `parentId` attribute to the value of the target realm’s ID.
   
   The newly added key is now the active key, as revealed by `kcadm.sh get keys -r demorealm`.

#### Adding new realm keys from a Java Key Store file

1. Add a new key provider to add a new key pair pre-prepared as a JKS file.
   
   For example, on:
   
   - Linux:
     
     ```
     $ kcadm.sh create components -r demorealm -s name=java-keystore -s providerId=java-keystore -s providerType=org.keycloak.keys.KeyProvider -s parentId=959844c1-d149-41d7-8359-6aa527fca0b0 -s 'config.priority=["101"]' -s 'config.enabled=["true"]' -s 'config.active=["true"]' -s 'config.keystore=["/opt/keycloak/keystore.jks"]' -s 'config.keystorePassword=["secret"]' -s 'config.keyPassword=["secret"]' -s 'config.keyAlias=["localhost"]'
     ```
   - Windows:
     
     ```
     c:\> kcadm create components -r demorealm -s name=java-keystore -s providerId=java-keystore -s providerType=org.keycloak.keys.KeyProvider -s parentId=959844c1-d149-41d7-8359-6aa527fca0b0 -s "config.priority=[\"101\"]" -s "config.enabled=[\"true\"]" -s "config.active=[\"true\"]" -s "config.keystore=[\"/opt/keycloak/keystore.jks\"]" -s "config.keystorePassword=[\"secret\"]" -s "config.keyPassword=[\"secret\"]" -s "config.keyAlias=[\"localhost\"]"
     ```
2. Ensure you change the attribute values for `keystore`, `keystorePassword`, `keyPassword`, and `alias` to match your specific keystore.
3. Set the `parentId` attribute to the value of the target realm’s ID.

#### Making the key passive or disabling the key

1. Identify the key you want to make passive.
   
   ```
   $ kcadm.sh get keys -r demorealm
   ```
2. Use the key’s `providerId` attribute to construct an endpoint URI, such as `components/PROVIDER_ID`.
3. Perform an `update`.
   
   For example:
   
   - Linux:
     
     ```
     $ kcadm.sh update components/PROVIDER_ID -r demorealm -s 'config.active=["false"]'
     ```
   - Windows:
     
     ```
     c:\> kcadm update components/PROVIDER_ID -r demorealm -s "config.active=[\"false\"]"
     ```
     
     You can update other key attributes:
   - Set a new `enabled` value to disable the key, for example, `config.enabled=["false"]`.
   - Set a new `priority` value to change the key’s priority, for example, `config.priority=["110"]`.

#### Deleting an old key

1. Ensure the key you are deleting is inactive and you have disabled it. This action is to prevent existing tokens held by applications and users from failing.
2. Identify the key to delete.
   
   ```
   $ kcadm.sh get keys -r demorealm
   ```
3. Use the `providerId` of the key to perform the delete.
   
   ```
   $ kcadm.sh delete components/PROVIDER_ID -r demorealm
   ```

#### Configuring event logging for a realm

Use the `update` command on the `events/config` endpoint.

The `eventsListeners` attribute contains a list of EventListenerProviderFactory IDs, specifying all event listeners that receive events. Attributes are available that control built-in event storage, so you can query past events using the Admin REST API. Keycloak has separate control over the logging of service calls (`eventsEnabled`) and the auditing events triggered by the Admin Console or Admin REST API (`adminEventsEnabled`). You can set up the `eventsExpiration` event to expire to prevent your database from filling. Keycloak sets `eventsExpiration` to time-to-live expressed in seconds.

You can set up a built-in event listener that receives all events and logs the events through JBoss-logging. Using the `org.keycloak.events` logger, Keycloak logs error events as `WARN` and other events as `DEBUG`.

For example:

- Linux:
  
  ```
  $ kcadm.sh update events/config -r demorealm -s 'eventsListeners=["jboss-logging"]'
  ```
- Windows:
  
  ```
  c:\> kcadm update events/config -r demorealm -s "eventsListeners=[\"jboss-logging\"]"
  ```

For example:

You can turn on storage for all available ERROR events, not including auditing events, for two days so you can retrieve the events through Admin REST.

- Linux:
  
  ```
  $ kcadm.sh update events/config -r demorealm -s eventsEnabled=true -s 'enabledEventTypes=["LOGIN_ERROR","REGISTER_ERROR","LOGOUT_ERROR","CODE_TO_TOKEN_ERROR","CLIENT_LOGIN_ERROR","FEDERATED_IDENTITY_LINK_ERROR","REMOVE_FEDERATED_IDENTITY_ERROR","UPDATE_EMAIL_ERROR","UPDATE_PROFILE_ERROR","UPDATE_PASSWORD_ERROR","UPDATE_TOTP_ERROR","UPDATE_CREDENTIAL_ERROR","VERIFY_EMAIL_ERROR","REMOVE_TOTP_ERROR","REMOVE_CREDENTIAL_ERROR","SEND_VERIFY_EMAIL_ERROR","SEND_RESET_PASSWORD_ERROR","SEND_IDENTITY_PROVIDER_LINK_ERROR","RESET_PASSWORD_ERROR","IDENTITY_PROVIDER_FIRST_LOGIN_ERROR","IDENTITY_PROVIDER_POST_LOGIN_ERROR","CUSTOM_REQUIRED_ACTION_ERROR","EXECUTE_ACTIONS_ERROR","CLIENT_REGISTER_ERROR","CLIENT_UPDATE_ERROR","CLIENT_DELETE_ERROR"]' -s eventsExpiration=172800
  ```
- Windows:
  
  ```
  c:\> kcadm update events/config -r demorealm -s eventsEnabled=true -s "enabledEventTypes=[\"LOGIN_ERROR\",\"REGISTER_ERROR\",\"LOGOUT_ERROR\",\"CODE_TO_TOKEN_ERROR\",\"CLIENT_LOGIN_ERROR\",\"FEDERATED_IDENTITY_LINK_ERROR\",\"REMOVE_FEDERATED_IDENTITY_ERROR\",\"UPDATE_EMAIL_ERROR\",\"UPDATE_PROFILE_ERROR\",\"UPDATE_PASSWORD_ERROR\",\"UPDATE_TOTP_ERROR\",\"UPDATE_CREDENTIAL_ERROR\",\"VERIFY_EMAIL_ERROR\",\"REMOVE_TOTP_ERROR\",\"REMOVE_CREDENTIAL_ERROR\",\"SEND_VERIFY_EMAIL_ERROR\",\"SEND_RESET_PASSWORD_ERROR\",\"SEND_IDENTITY_PROVIDER_LINK_ERROR\",\"RESET_PASSWORD_ERROR\",\"IDENTITY_PROVIDER_FIRST_LOGIN_ERROR\",\"IDENTITY_PROVIDER_POST_LOGIN_ERROR\",\"CUSTOM_REQUIRED_ACTION_ERROR\",\"EXECUTE_ACTIONS_ERROR\",\"CLIENT_REGISTER_ERROR\",\"CLIENT_UPDATE_ERROR\",\"CLIENT_DELETE_ERROR\"]" -s eventsExpiration=172800
  ```

You can reset stored event types to **all available event types**. Setting the value to an empty list is the same as enumerating all.

```
$ kcadm.sh update events/config -r demorealm -s enabledEventTypes=[]
```

You can enable storage of auditing events.

```
$ kcadm.sh update events/config -r demorealm -s adminEventsEnabled=true -s adminEventsDetailsEnabled=true
```

You can get the last 100 events. The events are ordered from newest to oldest.

```
$ kcadm.sh get events --offset 0 --limit 100
```

You can delete all saved events.

```
$ kcadm delete events
```

#### Flushing the caches

1. Use the `create` command with one of these endpoints to clear caches:
   
   - `clear-realm-cache`
   - `clear-user-cache`
   - `clear-keys-cache`
2. Set `realm` to the same value as the target realm.
   
   For example:
   
   ```
   $ kcadm.sh create clear-realm-cache -r demorealm -s realm=demorealm
   $ kcadm.sh create clear-user-cache -r demorealm -s realm=demorealm
   $ kcadm.sh create clear-keys-cache -r demorealm -s realm=demorealm
   ```

#### Importing a realm from exported .json file

1. Use the `create` command on the `partialImport` endpoint.
2. Set `ifResourceExists` to `FAIL`, `SKIP`, or `OVERWRITE`.
3. Use `-f` to submit the exported realm `.json` file.
   
   For example:
   
   ```
   $ kcadm.sh create partialImport -r demorealm2 -s ifResourceExists=FAIL -o -f demorealm.json
   ```
   
   If the realm does not yet exist, create it first.
   
   For example:
   
   ```
   $ kcadm.sh create realms -s realm=demorealm2 -s enabled=true
   ```

### [](#role-operations)Role operations

#### Creating a realm role

Use the `roles` endpoint to create a realm role.

```
$ kcadm.sh create roles -r demorealm -s name=user -s 'description=Regular user with a limited set of permissions'
```

#### Creating a client role

1. Identify the client.
2. Use the `get` command to list the available clients.
   
   ```
   $ kcadm.sh get clients -r demorealm --fields id,clientId
   ```
3. Create a new role by using the `clientId` attribute to construct an endpoint URI, such as `clients/ID/roles`.
   
   For example:
   
   ```
   $ kcadm.sh create clients/a95b6af3-0bdc-4878-ae2e-6d61a4eca9a0/roles -r demorealm -s name=editor -s 'description=Editor can edit, and publish any article'
   ```

#### Listing realm roles

Use the `get` command on the `roles` endpoint to list existing realm roles.

```
$ kcadm.sh get roles -r demorealm
```

You can use the `get-roles` command also.

```
$ kcadm.sh get-roles -r demorealm
```

#### Listing client roles

Keycloak has a dedicated `get-roles` command to simplify the listing of realm and client roles. The command is an extension of the `get` command and behaves the same as the `get` command but with additional semantics for listing roles.

Use the `get-roles` command by passing it the clientId (`--cclientid`) option or the `id` (`--cid`) option to identify the client to list client roles.

For example:

```
$ kcadm.sh get-roles -r demorealm --cclientid realm-management
```

#### Getting a specific realm role

Use the `get` command and the role `name` to construct an endpoint URI for a specific realm role, `roles/ROLE_NAME`, where `user` is the existing role’s name.

For example:

```
$ kcadm.sh get roles/user -r demorealm
```

You can use the `get-roles` command, passing it a role name (`--rolename` option) or ID (`--roleid` option).

For example:

```
$ kcadm.sh get-roles -r demorealm --rolename user
```

#### Getting a specific client role

Use the `get-roles` command, passing it the clientId attribute (`--cclientid` option) or ID attribute (`--cid` option) to identify the client, and pass the role name (`--rolename` option) or the role ID attribute (`--roleid`) to identify a specific client role.

For example:

```
$ kcadm.sh get-roles -r demorealm --cclientid realm-management --rolename manage-clients
```

#### Updating a realm role

Use the `update` command with the endpoint URI you used to get a specific realm role.

For example:

```
$ kcadm.sh update roles/user -r demorealm -s 'description=Role representing a regular user'
```

#### Updating a client role

Use the `update` command with the endpoint URI that you used to get a specific client role.

For example:

```
$ kcadm.sh update clients/a95b6af3-0bdc-4878-ae2e-6d61a4eca9a0/roles/editor -r demorealm -s 'description=User that can edit, and publish articles'
```

#### Deleting a realm role

Use the `delete` command with the endpoint URI that you used to get a specific realm role.

For example:

```
$ kcadm.sh delete roles/user -r demorealm
```

#### Deleting a client role

Use the `delete` command with the endpoint URI that you used to get a specific client role.

For example:

```
$ kcadm.sh delete clients/a95b6af3-0bdc-4878-ae2e-6d61a4eca9a0/roles/editor -r demorealm
```

#### Listing assigned, available, and effective realm roles for a composite role

Use the `get-roles` command to list assigned, available, and effective realm roles for a composite role.

1. To list **assigned** realm roles for the composite role, specify the target composite role by name (`--rname` option) or ID (`--rid` option).
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --rname testrole
   ```
2. Use the `--effective` option to list **effective** realm roles.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --rname testrole --effective
   ```
3. Use the `--available` option to list realm roles that you can add to the composite role.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --rname testrole --available
   ```

#### Listing assigned, available, and effective client roles for a composite role

Use the `get-roles` command to list assigned, available, and effective client roles for a composite role.

1. To list **assigned** client roles for the composite role, you can specify the target composite role by name (`--rname` option) or ID (`--rid` option) and client by the clientId attribute (`--cclientid` option) or ID (`--cid` option).
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --rname testrole --cclientid realm-management
   ```
2. Use the `--effective` option to list **effective** realm roles.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --rname testrole --cclientid realm-management --effective
   ```
3. Use the `--available` option to list realm roles that you can add to the target composite role.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --rname testrole --cclientid realm-management --available
   ```

#### Adding realm roles to a composite role

Keycloak provides an `add-roles` command for adding realm roles and client roles.

This example adds the `user` role to the composite role `testrole`.

```
$ kcadm.sh add-roles --rname testrole --rolename user -r demorealm
```

#### Removing realm roles from a composite role

Keycloak provides a `remove-roles` command for removing realm roles and client roles.

The following example removes the `user` role from the target composite role `testrole`.

```
$ kcadm.sh remove-roles --rname testrole --rolename user -r demorealm
```

#### Adding client roles to a realm role

Keycloak provides an `add-roles` command for adding realm roles and client roles.

The following example adds the roles defined on the client `realm-management`, `create-client`, and `view-users`, to the `testrole` composite role.

```
$ kcadm.sh add-roles -r demorealm --rname testrole --cclientid realm-management --rolename create-client --rolename view-users
```

#### Adding client roles to a client role

1. Determine the ID of the composite client role by using the `get-roles` command.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --cclientid test-client --rolename operations
   ```
2. Assume that a client exists with a clientId attribute named `test-client`, a client role named `support`, and a client role named `operations` which becomes a composite role that has an ID of "fc400897-ef6a-4e8c-872b-1581b7fa8a71".
3. Use the following example to add another role to the composite role.
   
   ```
   $ kcadm.sh add-roles -r demorealm --cclientid test-client --rid fc400897-ef6a-4e8c-872b-1581b7fa8a71 --rolename support
   ```
4. List the roles of a composite role by using the `get-roles --all` command.
   
   For example:
   
   ```
   $ kcadm.sh get-roles --rid fc400897-ef6a-4e8c-872b-1581b7fa8a71 --all
   ```

#### Removing client roles from a composite role

Use the `remove-roles` command to remove client roles from a composite role.

Use the following example to remove two roles defined on the client `realm-management`, the `create-client` role and the `view-users` role, from the `testrole` composite role.

```
$ kcadm.sh remove-roles -r demorealm --rname testrole --cclientid realm-management --rolename create-client --rolename view-users
```

#### Adding client roles to a group

Use the `add-roles` command to add realm roles and client roles.

The following example adds the roles defined on the client `realm-management`, `create-client` and `view-users`, to the `Group` group (`--gname` option). Alternatively, you can specify the group by ID (`--gid` option).

See [Group operations](#_group_operations) for more information.

```
$ kcadm.sh add-roles -r demorealm --gname Group --cclientid realm-management --rolename create-client --rolename view-users
```

#### Removing client roles from a group

Use the `remove-roles` command to remove client roles from a group.

The following example removes two roles defined on the client `realm-management`, `create-client` and `view-users`, from the `Group` group.

See [Group operations](#_group_operations) for more information.

```
$ kcadm.sh remove-roles -r demorealm --gname Group --cclientid realm-management --rolename create-client --rolename view-users
```

### [](#client-operations)Client operations

#### Creating a client

1. Run the `create` command on a `clients` endpoint to create a new client.
   
   For example:
   
   ```
   $ kcadm.sh create clients -r demorealm -s clientId=myapp -s enabled=true
   ```
2. Specify a secret if to set a secret for adapters to authenticate.
   
   For example:
   
   ```
   $ kcadm.sh create clients -r demorealm -s clientId=myapp -s enabled=true -s clientAuthenticatorType=client-secret -s secret=d0b8122f-8dfb-46b7-b68a-f5cc4e25d000
   ```

#### Listing clients

Use the `get` command on the `clients` endpoint to list clients.

This example filters the output to list only the `id` and `clientId` attributes:

```
$ kcadm.sh get clients -r demorealm --fields id,clientId
```

#### Getting a specific client

Use the client ID to construct an endpoint URI that targets a specific client, such as `clients/ID`.

For example:

```
$ kcadm.sh get clients/c7b8547f-e748-4333-95d0-410b76b3f4a3 -r demorealm
```

#### Getting the current secret for a specific client

Use the client ID to construct an endpoint URI, such as `clients/ID/client-secret`.

For example:

```
$ kcadm.sh get clients/$CID/client-secret -r demorealm
```

#### Generate a new secret for a specific client

Use the client ID to construct an endpoint URI, such as `clients/ID/client-secret`.

For example:

```
$ kcadm.sh create clients/$CID/client-secret -r demorealm
```

#### Updating the current secret for a specific client

Use the client ID to construct an endpoint URI, such as `clients/ID`.

For example:

```
$ kcadm.sh update clients/$CID -s "secret=newSecret" -r demorealm
```

#### Getting an adapter configuration file (keycloak.json) for a specific client

Use the client ID to construct an endpoint URI that targets a specific client, such as `clients/ID/installation/providers/keycloak-oidc-keycloak-json`.

For example:

```
$ kcadm.sh get clients/c7b8547f-e748-4333-95d0-410b76b3f4a3/installation/providers/keycloak-oidc-keycloak-json -r demorealm
```

#### Getting a WildFly subsystem adapter configuration for a specific client

Use the client ID to construct an endpoint URI that targets a specific client, such as `clients/ID/installation/providers/keycloak-oidc-jboss-subsystem`.

For example:

```
$ kcadm.sh get clients/c7b8547f-e748-4333-95d0-410b76b3f4a3/installation/providers/keycloak-oidc-jboss-subsystem -r demorealm
```

#### Getting a Docker-v2 example configuration for a specific client

Use the client ID to construct an endpoint URI that targets a specific client, such as `clients/ID/installation/providers/docker-v2-compose-yaml`.

The response is in `.zip` format.

For example:

```
$ kcadm.sh get http://localhost:8080/admin/realms/demorealm/clients/8f271c35-44e3-446f-8953-b0893810ebe7/installation/providers/docker-v2-compose-yaml -r demorealm > keycloak-docker-compose-yaml.zip
```

#### Updating a client

Use the `update` command with the same endpoint URI that you use to get a specific client.

For example:

- Linux:
  
  ```
  $ kcadm.sh update clients/c7b8547f-e748-4333-95d0-410b76b3f4a3 -r demorealm -s enabled=false -s publicClient=true -s 'redirectUris=["http://localhost:8080/myapp/*"]' -s baseUrl=http://localhost:8080/myapp -s adminUrl=http://localhost:8080/myapp
  ```
- Windows:
  
  ```
  c:\> kcadm update clients/c7b8547f-e748-4333-95d0-410b76b3f4a3 -r demorealm -s enabled=false -s publicClient=true -s "redirectUris=[\"http://localhost:8080/myapp/*\"]" -s baseUrl=http://localhost:8080/myapp -s adminUrl=http://localhost:8080/myapp
  ```

#### Deleting a client

Use the `delete` command with the same endpoint URI that you use to get a specific client.

For example:

```
$ kcadm.sh delete clients/c7b8547f-e748-4333-95d0-410b76b3f4a3 -r demorealm
```

#### Adding or removing roles for client’s service account

A client’s service account is a user account with username `service-account-CLIENT_ID`. You can perform the same user operations on this account as a regular account.

### [](#user-operations)User operations

#### Creating a user

Run the `create` command on the `users` endpoint to create a new user.

For example:

```
$ kcadm.sh create users -r demorealm -s username=testuser -s enabled=true
```

#### Listing users

Use the `users` endpoint to list users. The target user must change their password the next time they log in.

For example:

```
$ kcadm.sh get users -r demorealm --offset 0 --limit 1000
```

You can filter users by `username`, `firstName`, `lastName`, or `email`.

For example:

```
$ kcadm.sh get users -r demorealm -q q=email:google.com
$ kcadm.sh get users -r demorealm -q q=username:testuser
```

Filtering does not use exact matching. This example matches the value of the `username` attribute against the `*testuser*` pattern.

For clients, groups, and users you can filter across multiple attributes by specifying a more complex `q` query parameter. you may use something like -q q="field1:value1 field2:value2". Keycloak returns users that match the condition for all the attributes only.

#### Getting a specific user

Use the user ID to compose an endpoint URI, such as `users/USER_ID`.

For example:

```
$ kcadm.sh get users/0ba7a3fd-6fd8-48cd-a60b-2e8fd82d56e2 -r demorealm
```

#### Updating a user

Use the `update` command with the same endpoint URI that you use to get a specific user.

For example:

- Linux:
  
  ```
  $ kcadm.sh update users/0ba7a3fd-6fd8-48cd-a60b-2e8fd82d56e2 -r demorealm -s 'requiredActions=["VERIFY_EMAIL","UPDATE_PROFILE","CONFIGURE_TOTP","UPDATE_PASSWORD"]'
  ```
- Windows:
  
  ```
  c:\> kcadm update users/0ba7a3fd-6fd8-48cd-a60b-2e8fd82d56e2 -r demorealm -s "requiredActions=[\"VERIFY_EMAIL\",\"UPDATE_PROFILE\",\"CONFIGURE_TOTP\",\"UPDATE_PASSWORD\"]"
  ```

#### Deleting a user

Use the `delete` command with the same endpoint URI that you use to get a specific user.

For example:

```
$ kcadm.sh delete users/0ba7a3fd-6fd8-48cd-a60b-2e8fd82d56e2 -r demorealm
```

#### Resetting a user’s password

Use the dedicated `set-password` command to reset a user’s password.

For example:

```
$ kcadm.sh set-password -r demorealm --username testuser --new-password NEWPASSWORD --temporary
```

This command sets a temporary password for the user. The target user must change the password the next time they log in.

You can use `--userid` to specify the user by using the `id` attribute.

You can achieve the same result using the `update` command on an endpoint constructed from the one you used to get a specific user, such as `users/USER_ID/reset-password`.

For example:

```
$ kcadm.sh update users/0ba7a3fd-6fd8-48cd-a60b-2e8fd82d56e2/reset-password -r demorealm -s type=password -s value=NEWPASSWORD -s temporary=true -n
```

The `-n` parameter ensures that Keycloak performs the `PUT` command without performing a `GET` command before the `PUT` command. This is necessary because the `reset-password` endpoint does not support `GET`.

#### Listing assigned, available, and effective realm roles for a user

You can use a `get-roles` command to list assigned, available, and effective realm roles for a user.

1. Specify the target user by user name or ID to list the user’s **assigned** realm roles.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --uusername testuser
   ```
2. Use the `--effective` option to list **effective** realm roles.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --uusername testuser --effective
   ```
3. Use the `--available` option to list realm roles that you can add to a user.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --uusername testuser --available
   ```

#### Listing assigned, available, and effective client roles for a user

Use a `get-roles` command to list assigned, available, and effective client roles for a user.

1. Specify the target user by user name (`--uusername` option) or ID (`--uid` option) and client by a clientId attribute (`--cclientid` option) or an ID (`--cid` option) to list **assigned** client roles for the user.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --uusername testuser --cclientid realm-management
   ```
2. Use the `--effective` option to list **effective** realm roles.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --uusername testuser --cclientid realm-management --effective
   ```
3. Use the `--available` option to list realm roles that you can add to a user.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --uusername testuser --cclientid realm-management --available
   ```

#### Adding realm roles to a user

Use an `add-roles` command to add realm roles to a user.

Use the following example to add the `user` role to user `testuser`:

```
$ kcadm.sh add-roles --uusername testuser --rolename user -r demorealm
```

#### Removing realm roles from a user

Use a `remove-roles` command to remove realm roles from a user.

Use the following example to remove the `user` role from the user `testuser`:

```
$ kcadm.sh remove-roles --uusername testuser --rolename user -r demorealm
```

#### Adding client roles to a user

Use an `add-roles` command to add client roles to a user.

Use the following example to add two roles defined on the client `realm-management`, the `create-client` role and the `view-users` role, to the user `testuser`.

```
$ kcadm.sh add-roles -r demorealm --uusername testuser --cclientid realm-management --rolename create-client --rolename view-users
```

#### Removing client roles from a user

Use a `remove-roles` command to remove client roles from a user.

Use the following example to remove two roles defined on the realm-management client:

```
$ kcadm.sh remove-roles -r demorealm --uusername testuser --cclientid realm-management --rolename create-client --rolename view-users
```

#### Listing a user’s sessions

1. Identify the user’s ID,
2. Use the ID to compose an endpoint URI, such as `users/ID/sessions`.
3. Use the `get` command to retrieve a list of the user’s sessions.
   
   For example:
   
   ```
   $ kcadm.sh get users/6da5ab89-3397-4205-afaa-e201ff638f9e/sessions -r demorealm
   ```

#### Logging out a user from a specific session

1. Determine the session’s ID as described earlier.
2. Use the session’s ID to compose an endpoint URI, such as `sessions/ID`.
3. Use the `delete` command to invalidate the session.
   
   For example:
   
   ```
   $ kcadm.sh delete sessions/d0eaa7cc-8c5d-489d-811a-69d3c4ec84d1 -r demorealm
   ```

#### Logging out a user from all sessions

Use the user’s ID to construct an endpoint URI, such as `users/ID/logout`.

Use the `create` command to perform `POST` on that endpoint URI.

For example:

```
$ kcadm.sh create users/6da5ab89-3397-4205-afaa-e201ff638f9e/logout -r demorealm -s realm=demorealm -s user=6da5ab89-3397-4205-afaa-e201ff638f9e
```

### [](#_group_operations)Group operations

#### Creating a group

Use the `create` command on the `groups` endpoint to create a new group.

For example:

```
$ kcadm.sh create groups -r demorealm -s name=Group
```

#### Listing groups

Use the `get` command on the `groups` endpoint to list groups.

For example:

```
$ kcadm.sh get groups -r demorealm
```

#### Getting a specific group

Use the group’s ID to construct an endpoint URI, such as `groups/GROUP_ID`.

For example:

```
$ kcadm.sh get groups/51204821-0580-46db-8f2d-27106c6b5ded -r demorealm
```

#### Updating a group

Use the `update` command with the same endpoint URI that you use to get a specific group.

For example:

```
$ kcadm.sh update groups/51204821-0580-46db-8f2d-27106c6b5ded -s 'attributes.email=["group@example.com"]' -r demorealm
```

#### Deleting a group

Use the `delete` command with the same endpoint URI that you use to get a specific group.

For example:

```
$ kcadm.sh delete groups/51204821-0580-46db-8f2d-27106c6b5ded -r demorealm
```

#### Creating a subgroup

Find the ID of the parent group by listing groups. Use that ID to construct an endpoint URI, such as `groups/GROUP_ID/children`.

For example:

```
$ kcadm.sh create groups/51204821-0580-46db-8f2d-27106c6b5ded/children -r demorealm -s name=SubGroup
```

#### Moving a group under another group

1. Find the ID of an existing parent group and the ID of an existing child group.
2. Use the parent group’s ID to construct an endpoint URI, such as `groups/PARENT_GROUP_ID/children`.
3. Run the `create` command on this endpoint and pass the child group’s ID as a JSON body.

For example:

```
$ kcadm.sh create groups/51204821-0580-46db-8f2d-27106c6b5ded/children -r demorealm -s id=08d410c6-d585-4059-bb07-54dcb92c5094 -s name=SubGroup
```

#### Get groups for a specific user

Use a user’s ID to determine a user’s membership in groups to compose an endpoint URI, such as `users/USER_ID/groups`.

For example:

```
$ kcadm.sh get users/b544f379-5fc4-49e5-8a8d-5cfb71f46f53/groups -r demorealm
```

#### Adding a user to a group

Use the `update` command with an endpoint URI composed of a user’s ID and a group’s ID, such as `users/USER_ID/groups/GROUP_ID`, to add a user to a group.

For example:

```
$ kcadm.sh update users/b544f379-5fc4-49e5-8a8d-5cfb71f46f53/groups/ce01117a-7426-4670-a29a-5c118056fe20 -r demorealm -s realm=demorealm -s userId=b544f379-5fc4-49e5-8a8d-5cfb71f46f53 -s groupId=ce01117a-7426-4670-a29a-5c118056fe20 -n
```

#### Removing a user from a group

Use the `delete` command on the same endpoint URI you use for adding a user to a group, such as `users/USER_ID/groups/GROUP_ID`, to remove a user from a group.

For example:

```
$ kcadm.sh delete users/b544f379-5fc4-49e5-8a8d-5cfb71f46f53/groups/ce01117a-7426-4670-a29a-5c118056fe20 -r demorealm
```

#### Listing assigned, available, and effective realm roles for a group

Use a dedicated `get-roles` command to list assigned, available, and effective realm roles for a group.

1. Specify the target group by name (`--gname` option), path (`--gpath` option), or ID (`--gid` option) to list **assigned** realm roles for the group.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --gname Group
   ```
2. Use the `--effective` option to list **effective** realm roles.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --gname Group --effective
   ```
3. Use the `--available` option to list realm roles that you can add to the group.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --gname Group --available
   ```

#### Listing assigned, available, and effective client roles for a group

Use the `get-roles` command to list assigned, available, and effective client roles for a group.

1. Specify the target group by name (`--gname` option) or ID (`--gid` option),
2. Specify the client by the clientId attribute (`--cclientid` option) or ID (`--id` option) to list **assigned** client roles for the user.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --gname Group --cclientid realm-management
   ```
3. Use the `--effective` option to list **effective** realm roles.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --gname Group --cclientid realm-management --effective
   ```
4. Use the `--available` option to list realm roles that you can still add to the group.
   
   For example:
   
   ```
   $ kcadm.sh get-roles -r demorealm --gname Group --cclientid realm-management --available
   ```

### [](#identity-provider-operations)Identity provider operations

#### Listing available identity providers

Use the `serverinfo` endpoint to list available identity providers.

For example:

```
$ kcadm.sh get serverinfo -r demorealm --fields 'identityProviders(*)'
```

Keycloak processes the `serverinfo` endpoint similarly to the `realms` endpoint. Keycloak does not resolve the endpoint relative to a target realm because it exists outside any specific realm.

#### Listing configured identity providers

Use the `identity-provider/instances` endpoint.

For example:

```
$ kcadm.sh get identity-provider/instances -r demorealm --fields alias,providerId,enabled
```

#### Getting a specific configured identity provider

Use the identity provider’s `alias` attribute to construct an endpoint URI, such as `identity-provider/instances/ALIAS`, to get a specific identity provider.

For example:

```
$ kcadm.sh get identity-provider/instances/facebook -r demorealm
```

#### Removing a specific configured identity provider

Use the `delete` command with the same endpoint URI that you use to get a specific configured identity provider to remove a specific configured identity provider.

For example:

```
$ kcadm.sh delete identity-provider/instances/facebook -r demorealm
```

#### Configuring a Keycloak OpenID Connect identity provider

1. Use `keycloak-oidc` as the `providerId` when you create a new identity provider instance.
2. Provide the `config` attributes: `authorizationUrl`, `tokenUrl`, `clientId`, and `clientSecret`.
   
   For example:
   
   ```
   $ kcadm.sh create identity-provider/instances -r demorealm -s alias=keycloak-oidc -s providerId=keycloak-oidc -s enabled=true -s 'config.useJwksUrl="true"' -s config.authorizationUrl=http://localhost:8180/realms/demorealm/protocol/openid-connect/auth -s config.tokenUrl=http://localhost:8180/realms/demorealm/protocol/openid-connect/token -s config.clientId=demo-oidc-provider -s config.clientSecret=secret
   ```

#### Configuring an OpenID Connect identity provider

Configure the generic OpenID Connect provider the same way you configure the Keycloak OpenID Connect provider, except you set the `providerId` attribute value to `oidc`.

#### Configuring a SAML 2 identity provider

1. Use `saml` as the `providerId`.
2. Provide the `config` attributes: `singleSignOnServiceUrl`, `nameIDPolicyFormat`, and `signatureAlgorithm`.

For example:

```
$ kcadm.sh create identity-provider/instances -r demorealm -s alias=saml -s providerId=saml -s enabled=true -s 'config.useJwksUrl="true"' -s config.singleSignOnServiceUrl=http://localhost:8180/realms/saml-broker-realm/protocol/saml -s config.nameIDPolicyFormat=urn:oasis:names:tc:SAML:2.0:nameid-format:persistent -s config.signatureAlgorithm=RSA_SHA256
```

#### Configuring a Facebook identity provider

1. Use `facebook` as the `providerId`.
2. Provide the `config` attributes: `clientId` and `clientSecret`. You can find these attributes in the Facebook Developers application configuration page for your application. See the [Facebook identity broker](#_facebook) page for more information.
   
   For example:
   
   ```
   $ kcadm.sh create identity-provider/instances -r demorealm -s alias=facebook -s providerId=facebook -s enabled=true  -s 'config.useJwksUrl="true"' -s config.clientId=FACEBOOK_CLIENT_ID -s config.clientSecret=FACEBOOK_CLIENT_SECRET
   ```

#### Configuring a Google identity provider

1. Use `google` as the `providerId`.
2. Provide the `config` attributes: `clientId` and `clientSecret`. You can find these attributes in the Google Developers application configuration page for your application. See the [Google identity broker](#_google) page for more information.
   
   For example:
   
   ```
   $ kcadm.sh create identity-provider/instances -r demorealm -s alias=google -s providerId=google -s enabled=true  -s 'config.useJwksUrl="true"' -s config.clientId=GOOGLE_CLIENT_ID -s config.clientSecret=GOOGLE_CLIENT_SECRET
   ```

#### Configuring a Twitter identity provider

1. Use `twitter` as the `providerId`.
2. Provide the `config` attributes `clientId` and `clientSecret`. You can find these attributes in the Twitter Application Management application configuration page for your application. See the [Twitter identity broker](#_twitter) page for more information.
   
   For example:
   
   ```
   $ kcadm.sh create identity-provider/instances -r demorealm -s alias=google -s providerId=google -s enabled=true  -s 'config.useJwksUrl="true"' -s config.clientId=TWITTER_API_KEY -s config.clientSecret=TWITTER_API_SECRET
   ```

#### Configuring a GitHub identity provider

1. Use `github` as the `providerId`.
2. Provide the `config` attributes `clientId` and `clientSecret`. You can find these attributes in the GitHub Developer Application Settings page for your application. See the [GitHub identity broker](#_github) page for more information.
   
   For example:
   
   ```
   $ kcadm.sh create identity-provider/instances -r demorealm -s alias=github -s providerId=github -s enabled=true  -s 'config.useJwksUrl="true"' -s config.clientId=GITHUB_CLIENT_ID -s config.clientSecret=GITHUB_CLIENT_SECRET
   ```

#### Configuring a LinkedIn identity provider

1. Use `linkedin` as the `providerId`.
2. Provide the `config` attributes `clientId` and `clientSecret`. You can find these attributes in the LinkedIn Developer Console application page for your application. See the [LinkedIn identity broker](#_linkedin) page for more information.
   
   For example:
   
   ```
   $ kcadm.sh create identity-provider/instances -r demorealm -s alias=linkedin -s providerId=linkedin -s enabled=true  -s 'config.useJwksUrl="true"' -s config.clientId=LINKEDIN_CLIENT_ID -s config.clientSecret=LINKEDIN_CLIENT_SECRET
   ```

#### Configuring a Microsoft Live identity provider

1. Use `microsoft` as the `providerId`.
2. Provide the `config` attributes `clientId` and `clientSecret`. You can find these attributes in the Microsoft Application Registration Portal page for your application. See the [Microsoft identity broker](#_microsoft) page for more information.
   
   For example:
   
   ```
   $ kcadm.sh create identity-provider/instances -r demorealm -s alias=microsoft -s providerId=microsoft -s enabled=true  -s 'config.useJwksUrl="true"' -s config.clientId=MICROSOFT_APP_ID -s config.clientSecret=MICROSOFT_PASSWORD
   ```

#### Configuring a Stack Overflow identity provider

1. Use `stackoverflow` command as the `providerId`.
2. Provide the `config` attributes `clientId`, `clientSecret`, and `key`. You can find these attributes in the Stack Apps OAuth page for your application. See the [Stack Overflow identity broker](#_stackoverflow) page for more information.
   
   For example:
   
   ```
   $ kcadm.sh create identity-provider/instances -r demorealm -s alias=stackoverflow -s providerId=stackoverflow -s enabled=true  -s 'config.useJwksUrl="true"' -s config.clientId=STACKAPPS_CLIENT_ID -s config.clientSecret=STACKAPPS_CLIENT_SECRET -s config.key=STACKAPPS_KEY
   ```

### [](#storage-provider-operations)Storage provider operations

#### Configuring a Kerberos storage provider

1. Use the `create` command against the `components` endpoint.
2. Specify the realm id as a value of the `parentId` attribute.
3. Specify `kerberos` as the value of the `providerId` attribute, and `org.keycloak.storage.UserStorageProvider` as the value of the `providerType` attribute.
4. For example:
   
   ```
   $ kcadm.sh create components -r demorealm -s parentId=demorealmId -s id=demokerberos -s name=demokerberos -s providerId=kerberos -s providerType=org.keycloak.storage.UserStorageProvider -s 'config.priority=["0"]' -s 'config.debug=["false"]' -s 'config.allowPasswordAuthentication=["true"]' -s 'config.editMode=["UNSYNCED"]' -s 'config.updateProfileFirstLogin=["true"]' -s 'config.allowKerberosAuthentication=["true"]' -s 'config.kerberosRealm=["KEYCLOAK.ORG"]' -s 'config.keyTab=["http.keytab"]' -s 'config.serverPrincipal=["HTTP/localhost@KEYCLOAK.ORG"]' -s 'config.cachePolicy=["DEFAULT"]'
   ```

#### Configuring an LDAP user storage provider

1. Use the `create` command against the `components` endpoint.
2. Specify `ldap` as the value of the `providerId` attribute, and `org.keycloak.storage.UserStorageProvider` as the value of the `providerType` attribute.
3. Provide the realm ID as the value of the `parentId` attribute.
4. Use the following example to create a Kerberos-integrated LDAP provider.
   
   ```
   $ kcadm.sh create components -r demorealm -s name=kerberos-ldap-provider -s providerId=ldap -s providerType=org.keycloak.storage.UserStorageProvider -s parentId=3d9c572b-8f33-483f-98a6-8bb421667867  -s 'config.priority=["1"]' -s 'config.fullSyncPeriod=["-1"]' -s 'config.changedSyncPeriod=["-1"]' -s 'config.cachePolicy=["DEFAULT"]' -s config.evictionDay=[] -s config.evictionHour=[] -s config.evictionMinute=[] -s config.maxLifespan=[] -s 'config.batchSizeForSync=["1000"]' -s 'config.editMode=["WRITABLE"]' -s 'config.syncRegistrations=["false"]' -s 'config.vendor=["other"]' -s 'config.usernameLDAPAttribute=["uid"]' -s 'config.rdnLDAPAttribute=["uid"]' -s 'config.uuidLDAPAttribute=["entryUUID"]' -s 'config.userObjectClasses=["inetOrgPerson, organizationalPerson"]' -s 'config.connectionUrl=["ldap://localhost:10389"]'  -s 'config.usersDn=["ou=People,dc=keycloak,dc=org"]' -s 'config.authType=["simple"]' -s 'config.bindDn=["uid=admin,ou=system"]' -s 'config.bindCredential=["secret"]' -s 'config.searchScope=["1"]' -s 'config.useTruststoreSpi=["always"]' -s 'config.connectionPooling=["true"]' -s 'config.pagination=["true"]' -s 'config.allowKerberosAuthentication=["true"]' -s 'config.serverPrincipal=["HTTP/localhost@KEYCLOAK.ORG"]' -s 'config.keyTab=["http.keytab"]' -s 'config.kerberosRealm=["KEYCLOAK.ORG"]' -s 'config.debug=["true"]' -s 'config.useKerberosForPasswordAuthentication=["true"]'
   ```

#### Removing a user storage provider instance

1. Use the storage provider instance’s `id` attribute to compose an endpoint URI, such as `components/ID`.
2. Run the `delete` command against this endpoint.
   
   For example:
   
   ```
   $ kcadm.sh delete components/3d9c572b-8f33-483f-98a6-8bb421667867 -r demorealm
   ```

#### Triggering synchronization of all users for a specific user storage provider

1. Use the storage provider’s `id` attribute to compose an endpoint URI, such as `user-storage/ID_OF_USER_STORAGE_INSTANCE/sync`.
2. Add the `action=triggerFullSync` query parameter.
3. Run the `create` command.
   
   For example:
   
   ```
   $ kcadm.sh create user-storage/b7c63d02-b62a-4fc1-977c-947d6a09e1ea/sync?action=triggerFullSync
   ```

#### Triggering synchronization of changed users for a specific user storage provider

1. Use the storage provider’s `id` attribute to compose an endpoint URI, such as `user-storage/ID_OF_USER_STORAGE_INSTANCE/sync`.
2. Add the `action=triggerChangedUsersSync` query parameter.
3. Run the `create` command.
   
   For example:
   
   ```
   $ kcadm.sh create user-storage/b7c63d02-b62a-4fc1-977c-947d6a09e1ea/sync?action=triggerChangedUsersSync
   ```

#### Test LDAP user storage connectivity

1. Run the `get` command on the `testLDAPConnection` endpoint.
2. Provide query parameters `bindCredential`, `bindDn`, `connectionUrl`, and `useTruststoreSpi`.
3. Set the `action` query parameter to `testConnection`.
   
   For example:
   
   ```
   $ kcadm.sh create testLDAPConnection -s action=testConnection -s bindCredential=secret -s bindDn=uid=admin,ou=system -s connectionUrl=ldap://localhost:10389 -s useTruststoreSpi=always
   ```

#### Test LDAP user storage authentication

1. Run the `get` command on the `testLDAPConnection` endpoint.
2. Provide the query parameters `bindCredential`, `bindDn`, `connectionUrl`, and `useTruststoreSpi`.
3. Set the `action` query parameter to `testAuthentication`.
   
   For example:
   
   ```
   $ kcadm.sh create testLDAPConnection -s action=testAuthentication -s bindCredential=secret -s bindDn=uid=admin,ou=system -s connectionUrl=ldap://localhost:10389 -s useTruststoreSpi=always
   ```

### [](#adding-mappers)Adding mappers

#### Adding a hard-coded role LDAP mapper

1. Run the `create` command on the `components` endpoint.
2. Set the `providerType` attribute to `org.keycloak.storage.ldap.mappers.LDAPStorageMapper`.
3. Set the `parentId` attribute to the ID of the LDAP provider instance.
4. Set the `providerId` attribute to `hardcoded-ldap-role-mapper`. Ensure you provide a value of `role` configuration parameter.
   
   For example:
   
   ```
   $ kcadm.sh create components -r demorealm -s name=hardcoded-ldap-role-mapper -s providerId=hardcoded-ldap-role-mapper -s providerType=org.keycloak.storage.ldap.mappers.LDAPStorageMapper -s parentId=b7c63d02-b62a-4fc1-977c-947d6a09e1ea -s 'config.role=["realm-management.create-client"]'
   ```

#### Adding an MS Active Directory mapper

1. Run the `create` command on the `components` endpoint.
2. Set the `providerType` attribute to `org.keycloak.storage.ldap.mappers.LDAPStorageMapper`.
3. Set the `parentId` attribute to the ID of the LDAP provider instance.
4. Set the `providerId` attribute to `msad-user-account-control-mapper`.
   
   For example:
   
   ```
   $ kcadm.sh create components -r demorealm -s name=msad-user-account-control-mapper -s providerId=msad-user-account-control-mapper -s providerType=org.keycloak.storage.ldap.mappers.LDAPStorageMapper -s parentId=b7c63d02-b62a-4fc1-977c-947d6a09e1ea
   ```

#### Adding a user attribute LDAP mapper

1. Run the `create` command on the `components` endpoint.
2. Set the `providerType` attribute to `org.keycloak.storage.ldap.mappers.LDAPStorageMapper`.
3. Set the `parentId` attribute to the ID of the LDAP provider instance.
4. Set the `providerId` attribute to `user-attribute-ldap-mapper`.
   
   For example:
   
   ```
   $ kcadm.sh create components -r demorealm -s name=user-attribute-ldap-mapper -s providerId=user-attribute-ldap-mapper -s providerType=org.keycloak.storage.ldap.mappers.LDAPStorageMapper -s parentId=b7c63d02-b62a-4fc1-977c-947d6a09e1ea -s 'config."user.model.attribute"=["email"]' -s 'config."ldap.attribute"=["mail"]' -s 'config."read.only"=["false"]' -s 'config."always.read.value.from.ldap"=["false"]' -s 'config."is.mandatory.in.ldap"=["false"]'
   ```

#### Adding a group LDAP mapper

1. Run the `create` command on the `components` endpoint.
2. Set the `providerType` attribute to `org.keycloak.storage.ldap.mappers.LDAPStorageMapper`.
3. Set the `parentId` attribute to the ID of the LDAP provider instance.
4. Set the `providerId` attribute to `group-ldap-mapper`.
   
   For example:
   
   ```
   $ kcadm.sh create components -r demorealm -s name=group-ldap-mapper -s providerId=group-ldap-mapper -s providerType=org.keycloak.storage.ldap.mappers.LDAPStorageMapper -s parentId=b7c63d02-b62a-4fc1-977c-947d6a09e1ea -s 'config."groups.dn"=[]' -s 'config."group.name.ldap.attribute"=["cn"]' -s 'config."group.object.classes"=["groupOfNames"]' -s 'config."preserve.group.inheritance"=["true"]' -s 'config."membership.ldap.attribute"=["member"]' -s 'config."membership.attribute.type"=["DN"]' -s 'config."groups.ldap.filter"=[]' -s 'config.mode=["LDAP_ONLY"]' -s 'config."user.roles.retrieve.strategy"=["LOAD_GROUPS_BY_MEMBER_ATTRIBUTE"]' -s 'config."mapped.group.attributes"=["admins-group"]' -s 'config."drop.non.existing.groups.during.sync"=["false"]' -s 'config.roles=["admins"]' -s 'config.groups=["admins-group"]' -s 'config.group=[]' -s 'config.preserve=["true"]' -s 'config.membership=["member"]'
   ```

#### Adding a full name LDAP mapper

1. Run the `create` command on the `components` endpoint.
2. Set the `providerType` attribute to `org.keycloak.storage.ldap.mappers.LDAPStorageMapper`.
3. Set the `parentId` attribute to the ID of the LDAP provider instance.
4. Set the `providerId` attribute to `full-name-ldap-mapper`.
   
   For example:
   
   ```
   $ kcadm.sh create components -r demorealm -s name=full-name-ldap-mapper -s providerId=full-name-ldap-mapper -s providerType=org.keycloak.storage.ldap.mappers.LDAPStorageMapper -s parentId=b7c63d02-b62a-4fc1-977c-947d6a09e1ea -s 'config."ldap.full.name.attribute"=["cn"]' -s 'config."read.only"=["false"]' -s 'config."write.only"=["true"]'
   ```

### [](#authentication-operations)Authentication operations

#### Setting a password policy

1. Set the realm’s `passwordPolicy` attribute to an enumeration expression that includes the specific policy provider ID and optional configuration.
2. Use the following example to set a password policy to default values. The default values include:
   
   - 210,000 hashing iterations
   - at least one special character
   - at least one uppercase character
   - at least one digit character
   - not be equal to a user’s `username`
   - be at least eight characters long
     
     ```
     $ kcadm.sh update realms/demorealm -s 'passwordPolicy="hashIterations and specialChars and upperCase and digits and notUsername and length"'
     ```
3. To use values different from defaults, pass the configuration in brackets.
4. Use the following example to set a password policy to:
   
   - 300,000 hash iterations
   - at least two special characters
   - at least two uppercase characters
   - at least two lowercase characters
   - at least two digits
   - be at least nine characters long
   - not be equal to a user’s `username`
   - not repeat for at least four changes back
     
     ```
     $ kcadm.sh update realms/demorealm -s 'passwordPolicy="hashIterations(300000) and specialChars(2) and upperCase(2) and lowerCase(2) and digits(2) and length(9) and notUsername and passwordHistory(4)"'
     ```

#### Obtaining the current password policy

You can get the current realm configuration by filtering all output except for the `passwordPolicy` attribute.

For example, display `passwordPolicy` for `demorealm`.

```
$ kcadm.sh get realms/demorealm --fields passwordPolicy
```

#### Listing authentication flows

Run the `get` command on the `authentication/flows` endpoint.

For example:

```
$ kcadm.sh get authentication/flows -r demorealm
```

#### Getting a specific authentication flow

Run the `get` command on the `authentication/flows/FLOW_ID` endpoint.

For example:

```
$ kcadm.sh get authentication/flows/febfd772-e1a1-42fb-b8ae-00c0566fafb8 -r demorealm
```

#### Listing executions for a flow

Run the `get` command on the `authentication/flows/FLOW_ALIAS/executions` endpoint.

For example:

```
$ kcadm.sh get authentication/flows/Copy%20of%20browser/executions -r demorealm
```

#### Adding configuration to an execution

1. Get execution for a flow.
2. Note the ID of the flow.
3. Run the `create` command on the `authentication/executions/{executionId}/config` endpoint.

For example:

```
$ kcadm.sh create "authentication/executions/a3147129-c402-4760-86d9-3f2345e401c7/config" -r demorealm -b '{"config":{"x509-cert-auth.mapping-source-selection":"Match SubjectDN using regular expression","x509-cert-auth.regular-expression":"(.*?)(?:$)","x509-cert-auth.mapper-selection":"Custom Attribute Mapper","x509-cert-auth.mapper-selection.user-attribute-name":"usercertificate","x509-cert-auth.crl-checking-enabled":"","x509-cert-auth.crldp-checking-enabled":false,"x509-cert-auth.crl-relative-path":"crl.pem","x509-cert-auth.ocsp-checking-enabled":"","x509-cert-auth.ocsp-responder-uri":"","x509-cert-auth.keyusage":"","x509-cert-auth.extendedkeyusage":"","x509-cert-auth.confirmation-page-disallowed":""},"alias":"my_otp_config"}'
```

#### Getting configuration for an execution

1. Get execution for a flow.
2. Note its `authenticationConfig` attribute, which contains the config ID.
3. Run the `get` command on the `authentication/config/ID` endpoint.

For example:

```
$ kcadm get "authentication/config/dd91611a-d25c-421a-87e2-227c18421833" -r demorealm
```

#### Updating configuration for an execution

1. Get the execution for the flow.
2. Get the flow’s `authenticationConfig` attribute.
3. Note the config ID from the attribute.
4. Run the `update` command on the `authentication/config/ID` endpoint.

For example:

```
$ kcadm update "authentication/config/dd91611a-d25c-421a-87e2-227c18421833" -r demorealm -b '{"id":"dd91611a-d25c-421a-87e2-227c18421833","alias":"my_otp_config","config":{"x509-cert-auth.extendedkeyusage":"","x509-cert-auth.mapper-selection.user-attribute-name":"usercertificate","x509-cert-auth.ocsp-responder-uri":"","x509-cert-auth.regular-expression":"(.*?)(?:$)","x509-cert-auth.crl-checking-enabled":"true","x509-cert-auth.confirmation-page-disallowed":"","x509-cert-auth.keyusage":"","x509-cert-auth.mapper-selection":"Custom Attribute Mapper","x509-cert-auth.crl-relative-path":"crl.pem","x509-cert-auth.crldp-checking-enabled":"false","x509-cert-auth.mapping-source-selection":"Match SubjectDN using regular expression","x509-cert-auth.ocsp-checking-enabled":""}}'
```

#### Deleting configuration for an execution

1. Get execution for a flow.
2. Get the flows `authenticationConfig` attribute.
3. Note the config ID from the attribute.
4. Run the `delete` command on the `authentication/config/ID` endpoint.

For example:

```
$ kcadm delete "authentication/config/dd91611a-d25c-421a-87e2-227c18421833" -r demorealm
```

Last updated 2026-06-04 16:45:26 UTC