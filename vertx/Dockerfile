FROM azul/zulu-openjdk-alpine:21-latest

COPY ./build/libs/crate-cyclone-1.0.0-fat.jar /usr/app/
COPY ./.env /usr/app/

WORKDIR /usr/app

CMD ["java", "-jar", "crate-cyclone-1.0.0-fat.jar"]
