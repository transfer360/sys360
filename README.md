### Transfer360 System Structures for Golang

----

**ntk_to_lease** 
>Send Notice to Keeper data to lease company pubsub topic for onward processing

---

### Releasing New Versions

Go modules require semantic versioning format: `vMAJOR.MINOR.PATCH`

```bash
git tag v1.0.10
git push origin v1.0.10
```

Then consumers can update with:
```bash
go get github.com/transfer360/sys360@latest
```