package crate.event;

import java.util.Date;

public class Event {
    private Long id;
    private Long relationId;
    private Date time;
    private String detail;
    private String state;

    public Event(Builder builder) {

    }

    public static class Builder {
        private Long id;
        private Long relationId;
        private Date time;
        private String detail;
        private String state;

        public Builder id(Long id) {
            this.id = id;
            return this;
        }

        public Builder relationId(Long relationId) {
            this.relationId = relationId;
            return this;
        }

        public Builder time(Date time) {
            this.time = time;
            return this;
        }

        public Builder detail(String detail) {
            this.detail = detail;
            return this;
        }

        public Builder state(String state) {
            this.state = state;
            return this;
        }

        public Event build() {
            return new Event(this);
        }
    }

    public Long getId() {
        return id;
    }

    public Long getRelationId() {
        return relationId;
    }

    public Date getTime() {
        return time;
    }

    public String getDetail() {
        return detail;
    }

    public String getState() {
        return state;
    }
}
