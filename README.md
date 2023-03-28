aria2c --seed-time=0 https://cdimage.ubuntu.com/lubuntu/releases/20.04.2/release/lubuntu-20.04.2-desktop-amd64.iso.torrent
aria2c --seed-time=0 'magnet:?xt=urn:btih:015cc30f1e6b3d1e34069d33f09f8a2f25213495&dn=%E9%80%B1%E5%88%8A%E3%83%A4%E3%83%B3%E3%82%B0%E3%82%B8%E3%83%A3%E3%83%B3%E3%83%97%202022%E5%B9%B443%E5%8F%B7%20Weekly%20Young%20Jump%20No.43%20%5Bcomics888%5D&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce'

cd /Users/ohnishi/home/go/src/github.com/ohnishi/feed
go mod tidy


go run github.com/ohnishi/feed/cmd reset
go run github.com/ohnishi/feed/cmd feed
