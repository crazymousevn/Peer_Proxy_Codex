# Project Brief: Mạng Lưới Proxy Dân Dụng Lai (Hybrid P2P/Relay)

**Project Name:** (Gợi ý: PeerProxy)
**Version:** 1.0
**Date:** 2025-09-12

---

## 1. Executive Summary (Tóm tắt)

Dự án này nhằm mục đích xây dựng và kinh doanh một dịch vụ **mạng lưới proxy dân dụng (residential proxy) thế hệ mới**. Hệ thống sẽ tận dụng kết nối internet được chia sẻ từ một cộng đồng thành viên đông đảo để cung cấp địa chỉ IP đa dạng cho khách hàng.

Điểm khác biệt cốt lõi của dự án là áp dụng **kiến trúc lai (Hybrid Model)**, ưu tiên thiết lập kết nối **trực tiếp (P2P)** giữa khách hàng và thành viên để tối ưu tốc độ và giảm chi phí vận hành, đồng thời sử dụng **mô hình trung chuyển (Relay)** làm phương án dự phòng tin cậy, đảm bảo 100% kết nối thành công. Mục tiêu là tạo ra một dịch vụ proxy hiệu năng cao, khả năng mở rộng tốt và có lợi thế cạnh tranh về chi phí.

---

## 2. Problem Statement (Thách thức)

* **Nhu cầu Thị trường:** Khách hàng (doanh nghiệp, nhà nghiên cứu, cá nhân) có nhu cầu rất lớn về việc sử dụng proxy dân dụng để che giấu IP thật, phục vụ các tác vụ như thu thập dữ liệu web, nghiên cứu thị trường, quản lý tài khoản mạng xã hội mà không bị chặn.
* **Thách thức Kỹ thuật:** Nguồn cung cấp proxy (các thành viên) đều sử dụng mạng gia đình, có nghĩa là thiết bị của họ nằm sau **NAT (Network Address Translation)**. Điều này khiến việc thiết lập một kết nối proxy tiêu chuẩn (yêu cầu IP public và port mở) là **bất khả thi** nếu không có can thiệp kỹ thuật phức tạp từ phía thành viên.
* **Thách thức Vận hành:** Các mô hình proxy hoàn toàn dựa vào máy chủ trung chuyển (Relay) truyền thống gặp phải vấn đề lớn về **chi phí băng thông** và **khả năng mở rộng** khi lượng người dùng tăng cao, tạo ra điểm nghẽn cổ chai và làm tăng giá thành dịch vụ.

---

## 3. Proposed Solution (Giải pháp Đề xuất)

Chúng tôi sẽ xây dựng một hệ thống proxy thông minh dựa trên **Mô hình Lai (Hybrid Model)**, bao gồm các thành phần:

1.  **Ứng dụng Client cho Thành viên (Exit Node Client):** Một phần mềm zero-config, dễ dàng cài đặt trên máy tính của thành viên. Ứng dụng này sẽ chủ động kết nối ra ngoài và duy trì liên lạc với hệ thống trung tâm.
2.  **Hệ thống "Mai mối" Trung tâm (Signaling & Management System):**
    * **Signaling Server:** Không trung chuyển dữ liệu, chỉ làm nhiệm vụ "mai mối", trao đổi thông tin mạng để hai bên thiết lập kết nối trực tiếp.
    * **STUN/TURN Server:** Cung cấp các dịch vụ cần thiết để vượt NAT (NAT Traversal). Máy chủ TURN sẽ đóng vai trò dự phòng, trung chuyển dữ liệu khi P2P thất bại.
    * **Management Server:** Quản lý tài khoản, thanh toán, và theo dõi "sức khỏe" của mạng lưới.
3.  **Client-side Logic cho Khách hàng:** Cung cấp cho khách hàng một ứng dụng hoặc cơ chế để tham gia vào quá trình thiết lập kết nối P2P.

**Luồng hoạt động ưu tiên:**
* Hệ thống sẽ **luôn thử thiết lập kết nối P2P trực tiếp** giữa khách hàng và thành viên bằng công nghệ WebRTC/ICE.
* Nếu thành công (ước tính >80% trường hợp), lưu lượng truy cập sẽ đi thẳng, cho tốc độ tối đa và không tốn băng thông của hệ thống trung tâm.
* Nếu thất bại, hệ thống sẽ **tự động và liền mạch chuyển sang sử dụng máy chủ TURN** để trung chuyển dữ liệu, đảm bảo kết nối luôn được thiết lập.

---

## 4. Target Users (Đối tượng)

* **Khách hàng (End Users):**
    * Các doanh nghiệp cần thu thập dữ liệu web (web scraping) ở quy mô lớn.
    * Các chuyên gia marketing quản lý nhiều tài khoản mạng xã hội.
    * Các nhà nghiên cứu thị trường cần truy cập nội dung từ các vị trí địa lý khác nhau.
    * Người dùng cá nhân có nhu cầu cao về ẩn danh và bảo mật.
* **Thành viên (Exit Nodes):**
    * Người dùng internet thông thường trên toàn thế giới, có máy tính cá nhân và sẵn lòng chia sẻ một phần kết nối internet của mình (có thể để nhận lại một lợi ích nào đó như sử dụng dịch vụ miễn phí hoặc một khoản thù lao nhỏ).

---

## 5. Goals & Success Metrics (Mục tiêu & Thước đo)

* **Mục tiêu Kinh doanh:**
    * Ra mắt thành công phiên bản MVP trong vòng 6 tháng.
    * Đạt 1,000 thành viên và 100 khách hàng trả phí trong 3 tháng đầu sau ra mắt.
* **Mục tiêu Kỹ thuật & Sản phẩm:**
    * **Hiệu năng:** Tỷ lệ kết nối P2P thành công đạt trên 80%.
    * **Độ tin cậy:** Uptime của dịch vụ đạt 99.9%.
    * **Bảo mật:** Không có sự cố rò rỉ IP thật của khách hàng. Triển khai cơ chế theo dõi và chống lạm dụng hiệu quả.
    * **Trải nghiệm Người dùng:** Quá trình cài đặt cho thành viên chỉ mất dưới 2 phút. Khách hàng có thể tích hợp proxy vào các công cụ phổ biến một cách dễ dàng.

---

## 6. MVP Scope (Phạm vi cho Phiên bản Đầu tiên)

Phiên bản MVP sẽ tập trung vào việc chứng minh tính khả thi của mô hình lai.

* **Thành phần chính:**
    * Một **Signaling Server** cơ bản.
    * Một **STUN Server**.
    * Một **TURN Server** cơ bản (có thể dùng mã nguồn mở như `coturn`).
    * Một **Ứng dụng Client** cho thành viên (ban đầu có thể chỉ hỗ trợ Windows).
    * Một **Ứng dụng Client** đơn giản cho khách hàng để kiểm thử kết nối.
* **Tính năng cốt lõi:**
    * Khả năng thiết lập thành công kết nối P2P.
    * Khả năng tự động chuyển sang TURN khi P2P thất bại.
    * Xác thực người dùng và thành viên ở mức cơ bản.
    * Hệ thống gán proxy ngẫu nhiên (chưa cần lọc theo địa lý).

---

## 7. Post-MVP Vision (Tầm nhìn sau MVP)

* Mở rộng mạng lưới thành viên ra toàn cầu.
* Xây dựng dashboard quản lý hoàn chỉnh cho admin và khách hàng.
* Triển khai hệ thống thanh toán và các gói cước linh hoạt.
* Phát triển logic gán proxy thông minh (theo quốc gia, thành phố, ISP).
* Xây dựng hệ thống theo dõi "sức khỏe" IP và tự động luân chuyển các IP bị "cháy".
* Hỗ trợ nhiều nền tảng hơn cho ứng dụng client (macOS, Linux).

---

## 8. Constraints & Assumptions (Ràng buộc & Giả định)

* **Ràng buộc:**
    * Dự án có độ phức tạp kỹ thuật cao, đặc biệt ở phần mạng P2P.
    * Phụ thuộc vào việc xây dựng được một cộng đồng thành viên đủ lớn.
* **Giả định:**
    * Giả định rằng có thể thu hút đủ thành viên tham gia mạng lưới.
    * Giả định rằng tỷ lệ kết nối P2P thành công sẽ ở mức cao (>80%) để mô hình kinh doanh có hiệu quả về chi phí.

---

## 9. Risks & Open Questions (Rủi ro & Câu hỏi Mở)

* **Rủi ro Pháp lý:** Rủi ro lớn nhất là việc khách hàng lạm dụng hệ thống cho các hoạt động bất hợp pháp, làm ảnh hưởng đến các thành viên. Cần có chính sách pháp lý và kỹ thuật chặt chẽ để giảm thiểu.
* **Rủi ro Kỹ thuật:** Việc triển khai WebRTC/ICE ở quy mô lớn có thể gặp nhiều vấn đề không lường trước với các loại NAT khác nhau.
* **Rủi ro Cạnh tranh:** Bị các dịch vụ lớn (Google, Facebook) phát hiện và chặn hàng loạt thông qua phân tích hành vi.
* **Câu hỏi Mở:**
    * Mô hình đãi ngộ nào là hấp dẫn nhất để thu hút thành viên?
    * Làm thế nào để cân bằng giữa việc ghi log để chống lạm dụng và quyền riêng tư của khách hàng?

---

## 10. Next Steps (Bước tiếp theo)

* **Tạo Product Requirements Document (PRD):** Chi tiết hóa các yêu cầu chức năng và phi chức năng cho MVP dựa trên bản Project Brief này.
* **Thiết kế Kiến trúc Hệ thống Chi tiết:** Vẽ sơ đồ, lựa chọn công nghệ cụ thể cho từng thành phần.
* **Nghiên cứu Pháp lý:** Tham vấn về các điều khoản sử dụng dịch vụ cho cả hai bên.