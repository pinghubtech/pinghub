FROM scratch
ADD dist/pingbot_linux_arm64/pingbot /pingbot
EXPOSE 80
ENTRYPOINT ["/pinghub"]
