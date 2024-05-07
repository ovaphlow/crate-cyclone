package crate.setting;

import java.io.Serializable;

public class Setting implements Serializable {
    private Long id;
    private Long rootId;
    private Long parentId;
    private String tags;
    private String detail;
    private String state;

    public Setting(Long id, Long rootId, Long parentId, String tags, String detail, String state) {
        this.id = id;
        this.rootId = rootId;
        this.parentId = parentId;
        this.tags = tags;
        this.detail = detail;
        this.state = state;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Long getRootId() {
        return rootId;
    }

    public void setRootId(Long rootId) {
        this.rootId = rootId;
    }

    public Long getParentId() {
        return parentId;
    }

    public void setParentId(Long parentId) {
        this.parentId = parentId;
    }

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
        this.tags = tags;
    }

    public String getDetail() {
        return detail;
    }

    public void setDetail(String detail) {
        this.detail = detail;
    }

    public String getState() {
        return state;
    }

    public void setState(String state) {
        this.state = state;
    }
}
