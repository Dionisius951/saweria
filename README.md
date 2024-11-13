# Donate APP

#### Project ini adalah implementasi sederhana dari saweria dibuat dengan bahasa pemrograman golang dengan menggunakan implementasi dari UDP, TCP, dan WebSocket

## Cara install
* Clone repository ini

  https://github.com/Dionisius951/saweria.git

## Cara Menjalankan Server
1. Masuk ke directory server dengan cara "cd server"
2. Setelah masuk ke directory server buka terminal dan buka terminal lalu ketik "go run main.go" seperti pada contoh gambar dibawah ini
   ![{C0AC85B9-7397-48C1-A55A-5F21EBAF9E82}](https://github.com/user-attachments/assets/5d23713f-732e-409d-88e8-55c7672db4a9)
3. Jalankan file HTML dengan bantuan ekstensi dari vscode live server pada awal terbuka url hanya berupa localhost:PORT/index.html (PORT : port yang digunakan dalam live server anda)

   ![{A73921CB-63B9-4567-BD05-7894D8396378}](https://github.com/user-attachments/assets/d06db81c-c54d-4917-ba7b-2bcd899ab4df)

4. Tambahkan dibagian url seperti berikut localhost:PORT/index.html?username=[username] isi [username] dengan nama user yang diinginkan

   ![{36884278-4E83-4DFE-A853-079D31CB54DE}](https://github.com/user-attachments/assets/c9b887bc-6a46-436c-8164-0124c171dbfa)

5. Setelah itu tampilan akan seperti berikut dan halaman ini akan menampilkan setiap donasi yang masuk ke user tersebut

   ![{C297D35A-C19B-4078-A483-585CBAA8C4B0}](https://github.com/user-attachments/assets/e5cfb21c-4214-4f0c-9955-134f7d393179)

## Cara Menjalankan Client Untuk Melakukan Donasi
1. Buka terminal baru dan masuk ke directory udp dengan cara "cd client/udp"
2. Jalankan client dengan cara "go run udp_client.go [username] isi [username] dengan nama user yang diinginkan

   ![{C2804142-03BF-4293-8F79-548CAA1C01B8}](https://github.com/user-attachments/assets/35c04db6-6f6c-40f2-94c1-a78ad579cbcc)

## Cara Menjalankan Client Untuk Melakukan Topup
1. Buka terminal baru dan masuk ke directory tcp dengan cara "cd client/tcp"
2. Jalankan client dengan cara "go run tcp_client.go [username] isi [username] dengan nama user yang terdapat di udp_client

   ![{94E83D10-2C6A-434C-AA4B-8F2140BCB77C}](https://github.com/user-attachments/assets/912a0e12-09d2-42f1-b8bc-7e78a53bb829)


## Link Demo 

https://youtu.be/8vdmAFxRBnA



  



   
