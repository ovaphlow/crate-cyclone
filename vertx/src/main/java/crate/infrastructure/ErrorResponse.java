package crate.infrastructure;

import java.io.Serializable;

public class ErrorResponse implements Serializable {
    private String type;
    private Integer status;
    private String title;
    private String detail;
    private String instance;

    private ErrorResponse(Builder builder) {
        this.type = builder.type;
        this.status = builder.status;
        this.title = builder.title;
        this.detail = builder.detail;
        this.instance = builder.instance;
    }

    public static class Builder {
        private String type;
        private Integer status;
        private String title;
        private String detail;
        private String instance;

        public Builder type(String type) {
            this.type = type;
            return this;
        }

        public Builder status(Integer status) {
            this.status = status;
            return this;
        }

        public Builder title(String title) {
            this.title = title;
            return this;
        }

        public Builder detail(String detail) {
            this.detail = detail;
            return this;
        }

        public Builder instance(String instance) {
            this.instance = instance;
            return this;
        }

        public ErrorResponse build() {
            return new ErrorResponse(this);
        }
    }
}
