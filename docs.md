## CLIP ##
#### Pengertian ####
Clip adalah program yang berfungsi untuk menggali ringkasan informasi dai suatu web. 
Clip ditulis dengan menggunakan bahasa pemrograman `GO`. Cara kerja clip adalah
ketika ada request dari client untuk mencari _summary_ sebuah website, maka
clip akan membuka website tersebut dan membaca _summary_-nya berdasarkan :
1. OpenGraph Protocol >> http://ogp.me/
2. Twitter Cards >> https://dev.twitter.com/cards/overview
3. HTML Tag 

*Ketiganya dibaca menggunakan HTML Tokenizer di GO

#### Example JSON Result ####

```json
{
  "errors": null,
  "data": {
    "url": "https://www.youtube.com/watch?v=50efl4S8VQc",
    "title": "Sebangsa App (Versi Bahasa)",
    "description": "Sebangsa adalah platform sosial dan mobile baru yang dikembangkan dan dirancang untuk pengguna Indonesia. Gunakan Sebangsa untuk berhubungan dengan teman And...",
    "content_type": "html",
    "site_name": "YouTube",
    "image_thumb": "https://s.sebangsa.net/clip/i.ytimg.com/e95abe5367ab192d3c2f652aa4bff9c6.jpg",
    "image_width": "360",
    "image_height": "480",
    "favicon_url": "https://s.sebangsa.net/clip/s.ytimg.com/6a3f1f7deb777289735e2b44874f7bb0.png",
    "media_type": "video",
    "media_url": "https://s.sebangsa.net/clip/www.youtube.com/96a4a5b50aae6e025e035bba914c0161.50efl4S8VQc",
    "media_width": "1280",
    "media_height": "720",
    "published_date": "2014-07-25",
    "Client": null
  },
  "meta": null
}
```

## External Libary ##
- GIN https://github.com/gin-gonic/gin  (HTTP web framework)
- iconv https://github.com/qiniu/iconv (Convert string to requested character encoding)
- Logrus https://github.com/sirupsen/logrus (Structured, pluggable logging for Go)
- godotenv https://github.com/joho/godotenv (Loads environment variables from `.env`)
