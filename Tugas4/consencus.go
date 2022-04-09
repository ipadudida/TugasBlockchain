// Package consensus mengimplementasikan mesin kosensus etherium yang berbeda.
package consensus

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// ChainHeaderReader untuk mendefinisikan kumpulan kecil metorde yang diperlukan untuk mengakses blockchain lokal selama verivikasi header.

type ChainHeaderReader interface {
	// Config retrieves untuk mengkonfigurasi chain pada blockchain.
	Config() *params.ChainConfig

	// CurrentHeader untuk mengambil header saat ini pada chain lokal.
	CurrentHeader() *types.Header

	// GetHeader untuk mengambil header blok dari database dengan hash dan nomer.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber untuk mengambil header blok dari database dengan nomor.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash untuk mengambil header blok dari database dengan hashnya.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetTd retrieves untuk total kesulitan dari database dengan hash dan nomor.
	GetTd(hash common.Hash, number uint64) *big.Int
}

// ChainReader untuk mendefinisikan kumpulan kecil metode yang diperlukan untuk mengakses blockchain lokal selama verifikasi header dan/atau paman.
type ChainReader interface {
	ChainHeaderReader

	// GetBlock untuk mengambil blok dari database dengan hash dan nomor.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine adalah mesin konsensus agnostik algoritma.
type Engine interface {
	// Author berguna untuk mengambil alamat Ethereum dari akun yang mencetak blok yang diberikan, yang mungkin berbeda dari basis koin header jika mesin konsensus didasarkan pada tanda tangan.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader untuk memeriksa apakah header sesuai dengan aturan konsensus dari mesin yang diberikan. Memverifikasi segel dapat dilakukan secara opsional di sini, atau secara eksplisit melalui metode VerifySeal.
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders mirip dengan VerifyHeader, tetapi memverifikasi sekumpulan header secara bersamaan. Metode mengembalikan saluran keluar untuk membatalkan operasi dan saluran hasil untuk mengambil verifikasi asinkron (urutan adalah urutan irisan input).
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles untuk memverifikasi bahwa paman blok yang diberikan sesuai dengan aturan konsensus dari mesin yang diberikan.
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Prepare untuk menginisialisasi bidang konsensus dari header blok sesuai dengan aturan mesin tertentu. Perubahan dijalankan sebaris.
	Prepare(chain ChainHeaderReader, header *types.Header) error

	// Finalize untuk menjalankan modifikasi status pasca-transaksi (misalnya hadiah blok) tetapi tidak merakit blok. Catatan: Header blok dan database negara bagian mungkin diperbarui untuk mencerminkan aturan konsensus apa pun yang terjadi pada finalisasi (misalnya, hadiah blok).
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header)

	// FinalizeAndAssemble untuk menjalankan modifikasi status pasca-transaksi (misalnya hadiah blok) dan merakit blok terakhir. Catatan: Header blok dan database negara bagian mungkin diperbarui untuk mencerminkan aturan konsensus apa pun yang terjadi pada finalisasi (misalnya, hadiah blok).
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Seal untuk menghasilkan permintaan penyegelan baru untuk blok input yang diberikan dan mendorong hasilnya ke saluran yang diberikan. Catatan, metode ini segera kembali dan akan mengirimkan hasil async. Lebih dari satu hasil juga dapat dikembalikan tergantung pada algoritma konsensus.
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// SealHash untuk mengembalikan hash dari sebuah blok sebelum disegel.
	SealHash(header *types.Header) common.Hash

	// CalcDifficulty adalah algoritma penyesuaian kesulitan. Ini mengembalikan kesulitan yang seharusnya dimiliki blok baru.
	CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// APIs untuk mengembalikan API RPC yang disediakan mesin konsensus ini.
	APIs(chain ChainHeaderReader) []rpc.API

	// Close unttukmengakhiri utas latar belakang apa pun yang dikelola oleh mesin konsensus.
	Close() error
}

// PoW adalah mesin konsensus berdasarkan bukti kerja.
type PoW interface {
	Engine

	// Hashrate untuk mengembalikan hashrate penambangan saat ini dari mesin konsensus PoW.
	Hashrate() float64
}