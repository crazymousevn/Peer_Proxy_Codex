# PeerProxy Product Requirements Document (PRD)
**Version:** 1.0
**Date:** 2025-09-12

---

### **Phần 1: Mục tiêu và Bối cảnh**

#### **Mục tiêu**

* **Kinh doanh:**
    * Ra mắt phiên bản MVP thành công trong 6 tháng.
    * Đạt 1,000 thành viên (exit nodes) và 100 khách hàng trả phí trong 3 tháng đầu.
* **Sản phẩm & Kỹ thuật:**
    * Đạt tỷ lệ kết nối trực tiếp (P2P) thành công trên 80%.
    * Đảm bảo độ tin cậy của dịch vụ (uptime) đạt 99.5% cho MVP.
    * Bảo vệ tuyệt đối IP thật của khách hàng và xây dựng cơ chế chống lạm dụng hiệu quả.
    * Mang lại trải nghiệm người dùng đơn giản: thành viên cài đặt dưới 2 phút, khách hàng tích hợp dễ dàng.

#### **Bối cảnh**

Thị trường đang có nhu cầu lớn về proxy dân dụng để thu thập dữ liệu web, quản lý tài khoản mạng xã hội và nghiên cứu thị trường. Tuy nhiên, việc thiết lập kết nối proxy đến các máy tính cá nhân (nằm sau NAT) là một thách thức kỹ thuật lớn. Các giải pháp truyền thống dựa vào máy chủ trung chuyển (Relay) thì tốn kém và khó mở rộng.

Dự án PeerProxy ra đời để giải quyết vấn đề này bằng cách xây dựng một mạng lưới proxy dân dụng theo mô hình lai (Hybrid), ưu tiên kết nối P2P và tự động chuyển sang Relay khi cần, hứa hẹn tạo ra một dịch vụ hiệu năng cao với lợi thế cạnh tranh về giá.

#### **Nhật ký Thay đổi**

| Ngày       | Phiên bản | Mô tả                     | Tác giả    |
| :--------- | :-------- | :------------------------ | :-------- |
| 2025-09-12 | 1.0       | Hoàn tất phiên bản đầu tiên | John (PM) |

---

### **Phần 2: Yêu cầu (Requirements)**

#### **Yêu cầu Chức năng (Functional Requirements)**

* **FR1:** Ứng dụng phía thành viên (Exit Node Client) phải có khả năng tự động kết nối và duy trì liên lạc với Hệ thống Quản lý Trung tâm.
* **FR2:** Hệ thống phải triển khai một Signaling Server để trao đổi thông tin mạng (metadata) giữa khách hàng và thành viên.
* **FR3:** Hệ thống phải luôn ưu tiên thử thiết lập kết nối proxy trực tiếp (P2P) qua STUN.
* **FR4:** Khi kết nối P2P thất bại, hệ thống phải tự động chuyển sang sử dụng máy chủ TURN để trung chuyển (Relay) lưu lượng.
* **FR5:** Hệ thống phải có cơ chế xác thực cơ bản để định danh khách hàng và thành viên.
* **FR6:** Đối với MVP, hệ thống sẽ gán một thành viên (exit node) ngẫu nhiên cho khách hàng.
* **FR7:** Cung cấp một ứng dụng client cho thành viên (MVP hỗ trợ Windows) dạng dòng lệnh.
* **FR8:** Cung cấp một ứng dụng client cho khách hàng (MVP hỗ trợ Windows) có giao diện chẩn đoán.
* **FR9:** Hệ thống phải cung cấp cho khách hàng điểm cuối proxy theo giao thức SOCKS5 và HTTP/HTTPS.

#### **Yêu cầu Phi chức năng (Non-Functional Requirements)**

* **NFR1 (Hiệu năng):** Tỷ lệ kết nối P2P thành công phải đạt trên 80%.
* **NFR2 (Độ tin cậy):** Thời gian hoạt động (uptime) của dịch vụ phải đạt 99.5% cho MVP.
* **NFR3 (Bảo mật):** Hệ thống không được làm rò rỉ địa chỉ IP thật của khách hàng.
* **NFR4 (Trải nghiệm người dùng):** Quá trình cài đặt ứng dụng cho thành viên phải hoàn thành trong vòng dưới 2 phút.
* **NFR5 (Khả năng mở rộng):** Kiến trúc hệ thống phải được thiết kế để có thể mở rộng sau MVP.
* **NFR6 (Tính tương thích):** Khách hàng phải có khả năng tích hợp proxy vào các công cụ phổ biến.
* **NFR7.1 (Pháp lý):** Phải xây dựng Điều khoản Dịch vụ (Terms of Service) rõ ràng, nghiêm cấm mọi hành vi lạm dụng.
* **NFR7.2 (Chặn Port chủ động):** Ứng dụng client phía thành viên bắt buộc phải chặn lưu lượng trên port 25 (SMTP).
* **NFR7.3 (Ghi log minh bạch):** Hệ thống bắt buộc phải ghi lại metadata của kết nối (`ID khách hàng, ID thành viên, thời gian, IP/port đích`) và cam kết không ghi lại nội dung lưu lượng.

---

### **Phần 3: Mục tiêu Thiết kế Giao diện Người dùng (UI/UX)**

* **Tầm nhìn UX tổng thể (MVP):**
    * **Thành viên:** Không có giao diện, là một ứng dụng dòng lệnh chạy nền.
    * **Khách hàng:** Một công cụ tiện ích, đơn giản, tập trung vào việc kết nối, chẩn đoán và hiển thị trạng thái/lỗi một cách rõ ràng.
* **Màn hình Cốt lõi (MVP - Ứng dụng Khách hàng):**
    * Màn hình Kết nối với nút Bật/Tắt, tùy chọn "Buộc kết nối qua Relay", và hiển thị thông tin IP:Port, trạng thái, và log chẩn đoán.
* **Target Device and Platforms: Desktop Only (cho MVP)**
    * Cả hai ứng dụng sẽ chỉ hỗ trợ Windows trong phiên bản MVP.

---

### **Phần 4: Các Giả định Kỹ thuật (Technical Assumptions)**

* **Ngôn ngữ chính:** Toàn bộ logic nghiệp vụ của hệ thống sẽ được phát triển bằng **Go (Golang)**.
* **Client Framework:** Ứng dụng client có giao diện của khách hàng sẽ sử dụng framework **Wails**.
* **Cấu trúc Repository: Monorepo**.
* **Kiến trúc Dịch vụ: Monolith có cấu trúc** cho MVP.
* **Yêu cầu về Kiểm thử: Unit + Integration Tests** là bắt buộc.

---

### **Phần 5: Danh sách Epic**

* **Epic 1: Nền tảng và Kết nối Signaling**
* **Epic 2: Kết nối Proxy P2P & Xác thực Cơ bản**
* **Epic 3: Tăng cường Độ tin cậy & Khả năng Kiểm thử**
* **Epic 4: Sẵn sàng Ra mắt (Cơ chế Chống Lạm dụng Tối thiểu)**

---

### **Phần 6: Chi tiết Epic**

#### **Epic 1: Nền tảng và Kết nối Signaling**
* **Story 1.1:** Thiết lập Monorepo và Môi trường Phát triển.
* **Story 1.2:** Xây dựng Signaling Server Cơ bản.
* **Story 1.3:** Triển khai Trao đổi Thông điệp "Hello World".

#### **Epic 2: Kết nối Proxy P2P & Xác thực Cơ bản**
* **Story 2.1:** Xây dựng Hệ thống Xác thực Người dùng Cơ bản.
* **Story 2.2:** Tích hợp Thư viện WebRTC/ICE.
* **Story 2.3:** Triển khai Proxy SOCKS5 trên P2P Data Channel.
* **Story 2.4:** Triển khai Proxy HTTP/HTTPS trên P2P Data Channel.

#### **Epic 3: Tăng cường Độ tin cậy & Khả năng Kiểm thử**
* **Story 3.1:** Triển khai Kết nối qua TURN với Nút Chuyển đổi Thủ công.
* **Story 3.2:** Xây dựng Giao diện Người dùng Chẩn đoán cho Ứng dụng Khách hàng.

#### **Epic 4: Sẵn sàng Ra mắt (Cơ chế Chống Lạm dụng Tối thiểu)**
* **Story 4.1:** Triển khai Tính năng Chặn Port trong Ứng dụng Thành viên.
* **Story 4.2:** Xây dựng Hệ thống Ghi log Metadata Kết nối.

---

### **Phần 7: Bước tiếp theo (Next Steps)**

Tài liệu Yêu cầu Sản phẩm (PRD) này đã hoàn tất. Bước tiếp theo là chuyển giao cho Kiến trúc sư (Architect) để xây dựng Tài liệu Thiết kế Kiến trúc Hệ thống chi tiết.

**Gợi ý cho Kiến trúc sư:**
"Tài liệu PRD cho dự án PeerProxy đã hoàn tất. Vui lòng xem xét kỹ lưỡng để tạo ra một Tài liệu Kiến trúc Hệ thống chi tiết. Các giả định kỹ thuật chính cần tuân thủ bao gồm việc sử dụng Go cho tất cả các thành phần và cấu trúc monorepo. Bản kiến trúc cần chi tiết hóa việc triển khai Signaling Server, logic kết nối P2P/TURN, và các hệ thống backend cần thiết để hỗ trợ các story đã được định nghĩa."