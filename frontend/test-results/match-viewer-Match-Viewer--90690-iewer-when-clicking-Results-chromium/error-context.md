# Page snapshot

```yaml
- generic [active] [ref=e1]:
  - main [ref=e4]:
    - generic [ref=e5]:
      - generic [ref=e6]:
        - img "Langfuse Icon" [ref=e7]
        - heading "Sign in to your account" [level=2] [ref=e8]
      - generic [ref=e9]:
        - generic [ref=e12]:
          - generic [ref=e13]:
            - text: Email
            - textbox "Email" [ref=e14]:
              - /placeholder: jsdoe@example.com
          - generic [ref=e15]:
            - generic [ref=e16]:
              - text: Password
              - link "(forgot password?)" [ref=e17] [cursor=pointer]:
                - /url: /auth/reset-password
            - generic [ref=e18]:
              - textbox "Password (forgot password?)" [ref=e19]
              - button "Show password" [ref=e20] [cursor=pointer]:
                - img [ref=e21]
          - button "Sign in" [disabled] [ref=e24]
        - paragraph [ref=e25]:
          - text: No account yet?
          - link "Sign up" [ref=e26] [cursor=pointer]:
            - /url: /auth/sign-up?targetPath=%2Flogin
  - alert [ref=e27]: Sign in | Langfuse
```