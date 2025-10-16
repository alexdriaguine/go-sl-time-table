# sl api shit in go and maybe some framework later or idk 

## todo
- error handling in router.go for incorrect parsing etc
- set up an in memory thread safe cache (map with locks) wrapping sl-api, use that in router
- simple js front end, vanilla + gohtml (fetch on ssr also) -> select site -> save in cookie -> websocket?
- add sqlite for sites i guess. or also keep in memory, load at boot?
- add solidjs front end when above is done