cmake_minimum_required(VERSION 3.20)
project(esp32_prometheus_exporter)

# Find Go compiler
find_program(GO_EXECUTABLE go REQUIRED)
if(NOT GO_EXECUTABLE)
    message(FATAL_ERROR "Go compiler not found!")
endif()

# Set Go module path (adjust if needed)
set(GO_MODULE_PATH "${CMAKE_SOURCE_DIR}")

# Define the Go executable target
add_custom_target(build-go-binary
    COMMAND ${GO_EXECUTABLE} mod tidy
    COMMAND ${GO_EXECUTABLE} build -o ${CMAKE_BINARY_DIR}/esp32_prometheus_exporter cmd/app/main.go
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Building Go binary from cmd/app/main.go"
    VERBATIM
)

# Install binary
install(PROGRAMS ${CMAKE_BINARY_DIR}/esp32_prometheus_exporter
        DESTINATION bin
        COMPONENT runtime
        OPTIONAL
)

# Generate and install systemd service
set(SERVICE_FILE ${CMAKE_BINARY_DIR}/esp32_prometheus_exporter.service)
configure_file(esp32_prometheus_exporter.service.in ${SERVICE_FILE})
install(FILES ${SERVICE_FILE} DESTINATION /etc/systemd/system/)
