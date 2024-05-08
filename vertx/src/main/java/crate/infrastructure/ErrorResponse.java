package crate.infrastructure;

import java.io.Serializable;

public record ErrorResponse(
    String type,
    Integer status,
    String title,
    String detail,
    String instance
) implements Serializable {
    public ErrorResponse {
        if (null == type) {
            type = "abount:blank";
        }
    }
}
