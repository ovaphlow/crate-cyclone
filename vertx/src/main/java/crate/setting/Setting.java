package crate.setting;

import java.io.Serializable;
import java.time.OffsetDateTime;

public class Setting implements Serializable {
    private final Long id;
    private final Long rootId;
    private final Long parentId;
    private final String tags;
    private final String detail;
    private final OffsetDateTime createdAt;
    private final OffsetDateTime updatedAt;
    private final String state;

    private Setting(Builder builder) {
        this.id = builder.id;
        this.rootId = builder.rootId;
        this.parentId = builder.parentId;
        this.tags = builder.tags;
        this.detail = builder.detail;
        this.createdAt = builder.createdAt;
        this.updatedAt = builder.updatedAt;
        this.state = builder.state;
    }

    public Long getId() {
        return id;
    }

    public Long getRootId() {
        return rootId;
    }

    public Long getParentId() {
        return parentId;
    }

    public String getTags() {
        return tags;
    }

    public String getDetail() {
        return detail;
    }

    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    public String getState() {
        return state;
    }

    @Override
    public String toString() {
        return "Setting{" +
            "id=" + id +
            ", rootId=" + rootId +
            ", parentId=" + parentId +
            ", tags='" + tags + '\'' +
            ", detail='" + detail + '\'' +
            ", createdAt=" + createdAt +
            ", updatedAt=" + updatedAt +
            ", state='" + state + '\'' +
            '}';
    }

    public static class Builder {
        private Long id;
        private Long rootId;
        private Long parentId;
        private String tags;
        private String detail;
        private OffsetDateTime createdAt;
        private OffsetDateTime updatedAt;
        private String state;

        public Builder id(Long id) {
            this.id = id;
            return this;
        }

        public Builder rootId(Long rootId) {
            this.rootId = rootId;
            return this;
        }

        public Builder parentId(Long parentId) {
            this.parentId = parentId;
            return this;
        }

        public Builder tags(String tags) {
            this.tags = tags;
            return this;
        }

        public Builder detail(String detail) {
            this.detail = detail;
            return this;
        }

        public Builder createdAt(OffsetDateTime createdAt) {
            this.createdAt = createdAt;
            return this;
        }

        public Builder updatedAt(OffsetDateTime updatedAt) {
            this.updatedAt = updatedAt;
            return this;
        }

        public Builder state(String state) {
            this.state = state;
            return this;
        }

        public Setting build() {
            return new Setting(this);
        }
    }
}
