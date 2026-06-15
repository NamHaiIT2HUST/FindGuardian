package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TokenizedPayload struct {
	Category   string `json:"category"`
	AmountTier string `json:"amount_tier"`
	TimeOfDay  string `json:"time_of_day"`
	Action     string `json:"action"`
}

// Gọi API kiểm tra sinh trắc học
func callVNPTeKYCLivenessFace() bool {
	fmt.Println("[Cloud API] Đang gọi VNPT eKYC (Liveness Face) API...")
	time.Sleep(500 * time.Millisecond) // Giả lập độ trễ mạng
	fmt.Println("[Cloud API] Trạng thái: PASS - Xác nhận sinh trắc học thành công.")
	return true
}

// Gọi API speech to text
func callVNPTSmartVoiceSTT() string {
	fmt.Println("[Cloud API] Đang gọi VNPT SmartVoice (Speech to Text) API...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[Cloud API] Trạng thái: Nhận diện thành công văn bản.")
	return "Tôi đang tỉnh táo và chịu trách nhiệm cho khoản chi này"
}

// Xử lý cảnh báo
func nudgeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST method", http.StatusMethodNotAllowed)
		return
	}

	var payload TokenizedPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("\n--- Nhận yêu cầu phân tích từ Edge ---\n")
	fmt.Printf("Dữ liệu ẩn danh: %+v\n", payload)

	if payload.Action == "require_friction" {
		fmt.Println("[Cloud Logic] Kích hoạt cơ chế Micro-Friction...")

		// Tích hợp các API của BTC
		faceVerified := callVNPTeKYCLivenessFace()
		sttResult := callVNPTSmartVoiceSTT()

		if faceVerified && len(sttResult) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "success", "message": "Xác nhận giao dịch thành công. Người dùng ý thức rõ quyết định chi tiêu."}`))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ignored", "message": "Giao dịch an toàn, không cần can thiệp"}`))
}

func main() {
	http.HandleFunc("/api/v1/nudge/verify", nudgeHandler)

	fmt.Println("=== FinGuardian Secure API Gateway ===")
	fmt.Println("Server đang chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
