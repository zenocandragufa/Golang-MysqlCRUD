package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

//membuat type mahasiswa dengan struktur
type mahasiswa struct {
	Nim    string
	Nama   string
	Progdi string
	Smt    int
}

//membuat type response dengan struktur
type response struct {
	Status bool
	Pesan  string
	Data   []mahasiswa
}

//membuat fungsi koneksi dengan sql
//sintax -> sql.Open("mysql", "user:password@tcp(host:port)/nama_database")
//karena bawaan xampp password kosong jd dikosongkan aja
func koneksi() (*sql.DB, error) {
	db, salahe := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/cloud_udb")
	if salahe != nil {
		return nil, salahe
	}
	return db, nil
}

//fungsi tampil data
func tampil(pesane string) response {
	db, salahe := koneksi()
	if salahe != nil {

		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	defer db.Close()
	dataMhs, salahe := db.Query("select * from mahasiswa")
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	defer dataMhs.Close()
	var hasil []mahasiswa
	for dataMhs.Next() {
		var mhs = mahasiswa{}
		var salahe = dataMhs.Scan(&mhs.Nim, &mhs.Nama, &mhs.Progdi, &mhs.Smt)
		if salahe != nil {
			return response{
				Status: false,
				Pesan:  "Gagal Baca: " + salahe.Error(),
				Data:   []mahasiswa{},
			}
		}
		hasil = append(hasil, mhs)
	}
	salahe = dataMhs.Err()
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Kesalahan: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	return response{
		Status: true,
		Pesan:  pesane,
		Data:   hasil,
	}

}

//fungsi tampil data berdasarkan nim
func getMhs(nim string) response {
	db, salahe := koneksi()
	if salahe != nil {

		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	defer db.Close()

	dataMhs, salahe := db.Query("select * from mahasiswa where nim=?", nim)
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	defer dataMhs.Close()
	var hasil []mahasiswa

	for dataMhs.Next() {
		var mhs = mahasiswa{}
		var salahe = dataMhs.Scan(&mhs.Nim, &mhs.Nama, &mhs.Progdi, &mhs.Smt)
		if salahe != nil {
			return response{
				Status: false,
				Pesan:  "Gagal Baca: " + salahe.Error(),
				Data:   []mahasiswa{},
			}
		}
		hasil = append(hasil, mhs)
	}
	salahe = dataMhs.Err()

	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Kesalahan: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil Tampil",
		Data:   hasil,
	}

}

//fungsi tambah data
func tambah(nim string, nama string, progdi string, smt string) response {
	db, salahe := koneksi()
	if salahe != nil {

		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	defer db.Close()

	_, salahe = db.Exec("insert into mahasiswa values (?, ?, ?, ?)", nim, nama, progdi, smt)
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query Insert: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil Tambah",
		Data:   []mahasiswa{},
	}

}

//fungsi ubah data
func ubah(nim string, nama string, progdi string, smt string) response {
	db, salahe := koneksi()
	if salahe != nil {

		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	defer db.Close()

	_, salahe = db.Exec("update mahasiswa set nama=?, progdi=?, smt=? where nim=?", nama, progdi, smt, nim)
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query Update: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil Ubah",
		Data:   []mahasiswa{},
	}

}

//fungsi hapus data
func hapus(nim string) response {
	db, salahe := koneksi()
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	defer db.Close()

	_, salahe = db.Exec("delete from mahasiswa where nim=?", nim)
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query Delete: " + salahe.Error(),
			Data:   []mahasiswa{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil Hapus",
		Data:   []mahasiswa{},
	}

}
func kontroler(w http.ResponseWriter, r *http.Request) {
	var tampilHtml, salaheTampil = template.ParseFiles("template/tampil.html")
	if salaheTampil != nil {
		fmt.Println(salaheTampil.Error())
		return
	}
	var tambahHtml, salaheTambah = template.ParseFiles("template/tambah.html")
	if salaheTambah != nil {
		fmt.Println(salaheTambah.Error())
		return
	}
	var ubahHtml, salaheUbah = template.ParseFiles("template/ubah.html")
	if salaheUbah != nil {
		fmt.Println(salaheUbah.Error())
		return
	}
	var hapusHtml, salaheHapus = template.ParseFiles("template/hapus.html")
	if salaheHapus != nil {
		fmt.Println(salaheHapus.Error())
		return
	}
	switch r.Method {
	case "GET":
		aksi := r.URL.Query()["aksi"]
		if len(aksi) == 0 {
			tampilHtml.Execute(w, tampil("Berhasil Tampil"))

		} else if aksi[0] == "tambah" {
			tambahHtml.Execute(w, nil)

		} else if aksi[0] == "ubah" {
			nim := r.URL.Query()["nim"]
			ubahHtml.Execute(w, getMhs(nim[0]))

		} else if aksi[0] == "hapus" {
			nim := r.URL.Query()["nim"]
			hapusHtml.Execute(w, getMhs(nim[0]))

		} else {
			tampilHtml.Execute(w, tampil("Berhasil Tampil"))
		}

	case "POST":
		var salahe = r.ParseForm()
		if salahe != nil {
			fmt.Fprintln(w, "Kesalahan: ", salahe)
			return
		}

		var nim = r.FormValue("nim")
		var nama = r.FormValue("nama")
		var progdi = r.FormValue("progdi")
		var smt = r.FormValue("smt")

		var aksi = r.URL.Path
		if aksi == "/tambah" {
			var hasil = tambah(nim, nama, progdi, smt)
			tampilHtml.Execute(w, tampil(hasil.Pesan))
		} else if aksi == "/ubah" {
			var hasil = ubah(nim, nama, progdi, smt)
			tampilHtml.Execute(w, tampil(hasil.Pesan))
		} else if aksi == "/hapus" {
			var hasil = hapus(nim)
			tampilHtml.Execute(w, tampil(hasil.Pesan))
		} else {
			tampilHtml.Execute(w, tampil("Berhasil Tampil"))
		}

	default:
		fmt.Fprint(w, "Maaf. Method yang didukung hanya GET dan POST")
	}

}
func main() {
	http.HandleFunc("/", kontroler)
	fmt.Println("Server berjalan di Port 8080...")
	http.ListenAndServe(":8080", nil)
}
