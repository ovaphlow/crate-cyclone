package crate.infrastructure.model;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Builder
public class ErrorResponse {

    private String type;
    private Integer status;
    private String title;
    private String detail;
    private String instance;
}
