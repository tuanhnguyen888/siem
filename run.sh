#!/bin/bash

# Kiểm tra xem Docker Compose đã được cài đặt hay chưa
if [ ! -x "$(command -v docker-compose)" ]; then
    echo "Docker Compose không được cài đặt. Vui lòng cài đặt Docker Compose trước khi chạy script này."
    exit 1
fi

# Danh sách các thư mục project
project_directories=("kafkService" "mongodbService", "filter/logstash" ,  "collector", "correl", "go-api", "portal")

# Lặp qua từng thư mục và chạy Docker Compose
for project_dir in "${project_directories[@]}"; do
    cd "$project_dir" || exit 1

    echo "Bắt đầu build và chạy project trong thư mục: $project_dir"

    # Build và chạy Docker Compose trong từng thư mục project
    docker-compose build
    docker-compose up -d

    echo "Project trong thư mục $project_dir đã được build và chạy thành công."

    # Di chuyển trở lại thư mục gốc
    cd ..
done

echo "Tất cả các project đã được build và chạy thành công."
