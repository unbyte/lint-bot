# Just for fun

```yaml
auth:
  pat: ${github personal access token}
  secret: ${project webhook secret}
rules:
  - consume: log
    produce: text
    formatters:
      - unfold
```