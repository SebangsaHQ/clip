# Clip

Service yang dapat menggali informasi ringkas suatu website

## Cara Penggunaan ##
1. Pastikan Golang sudah terinstall pada PC
2. Clone Repo ini
3. Jalankan perintah `go get .` on pada repo yang telah di-clone
4. Lakikan pengaturan file .env *lihat caranya di bawah
5. jalankan program dengan menggunakan perintah `./run.sh serve`
6. coba lakukan POST request dengan `url` sebagai parameter. see this example :
   ```
   curl -X POST -F 'url=https://youtu.be/50efl4S8VQc' http://localhost:3001
   ```
7. Lihat hasilnya

## Mengkonfigurasi _environment_ : ##
- Buka contoh file `.env` yang bernama `.env.example`
- Edit file sesuai kebutuhan
- save as dengan nama ".env"

## Dokumentasi ##
[Lihat Dokumentasi](https://github.com/SebangsaHQ/clip/wiki)

## Lisensi ##

This package is licensed under MIT license. See `LICENSE.txt` for details.
