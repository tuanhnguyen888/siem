# siem
# Hướng dẫn sử dụng Script Build và Run Docker Compose

## Mục đích

Script này được tạo để tự động hóa quá trình build và chạy Docker Compose cho các project của bạn.

## Yêu Cầu

1. Docker Compose đã được cài đặt trên máy tính của bạn.

## Sử Dụng


1. Gán quyền thực thi cho script.

    ```bash
    chmod +x run.sh
    ```

2. Chạy script.

    ```bash
    ./run.sh
    ```

3. Script sẽ tự động build và chạy Docker Compose cho từng project trong danh sách các thư mục đã được định nghĩa trong script.

4. Khi quá trình hoàn thành, tất cả các project sẽ được build và chạy thành công.

Chú ý: Trước khi chạy script, đảm bảo rằng tệp `docker-compose.yml` của mỗi project được cấu hình đúng và không có lỗi xảy ra khi chạy Docker Compose.

