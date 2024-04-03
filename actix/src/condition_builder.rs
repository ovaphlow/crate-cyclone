pub fn equal_builder(equal: &[&str]) -> (Vec<String>, Vec<String>) {
    if equal.len() % 2 != 0 {
        print!("{}\n", "equal length must be even");
        return (Vec::new(), Vec::new());
    }
    let mut conditions = Vec::new();
    let mut params = Vec::new();
    for i in (0..equal.len()).step_by(2) {
        let c = format!("{} = ?", equal[i]);
        conditions.push(c);
        params.push(equal[i+1].to_string());
    }
    (conditions, params)
}

pub fn object_contain_builder(object_contain: &[&str]) -> (Vec<String>, Vec<String>) {
    if object_contain.len() % 3 != 0 {
        print!("{}\n", "object_contain length must be multiple of 3");
        return (Vec::new(), Vec::new())
    }
    let mut conditions = Vec::new();
    let mut params = Vec::new();
    for i in (0..object_contain.len()).step_by(3) {
        let c = format!("json_contains({}, json_object('{}', ?))", object_contain[i], object_contain[i+1]);
        conditions.push(c);
        params.push(object_contain[i+2].to_string());
    }
    (conditions, params)
}

pub fn array_contain_builder(array_contain: &[&str]) -> (Vec<String>, Vec<String>) {
    if array_contain.len() % 2 != 0 {
        print!("{}\n", "array_contain length must be even");
        return (Vec::new(), Vec::new())
    }
    let mut conditions = Vec::new();
    let mut params = Vec::new();
    for i in (0..array_contain.len()).step_by(2) {
        let c = format!("json_contains({}, json_array(?))", array_contain[i]);
        conditions.push(c);
        params.push(array_contain[i+1].to_string());
    }
    (conditions, params)
}

pub fn like_builder(like: &[&str]) -> (Vec<String>, Vec<String>) {
    if like.len() % 2 != 0 {
        print!("{}\n", "like length must be even");
        return (Vec::new(), Vec::new())
    }
    let mut conditions = Vec::new();
    let mut params = Vec::new();
    for i in (0..like.len()).step_by(2) {
        let c = format!("position(? in {})", like[i]);
        conditions.push(c);
        params.push(like[i+1].to_string());
    }
    (conditions, params)
}

pub fn object_like_builder(object_like: &[&str]) -> (Vec<String>, Vec<String>) {
    if object_like.len() % 3 != 0 {
        print!("{}\n", "object_like length must be multiple of 3");
        return (Vec::new(), Vec::new())
    }
    let mut conditions = Vec::new();
    let mut params = Vec::new();
    for i in (0..object_like.len()).step_by(3) {
        let c = format!("position(? in {}->>'$.{}')", object_like[i], object_like[i+1]);
        conditions.push(c);
        params.push(object_like[i+2].to_string());
    }
    (conditions, params)
}

pub fn in_builder(in_: &[&str]) -> (Vec<String>, Vec<String>) {
    if in_.len() < 2 {
        print!("{}\n", "in length must be even");
        return (Vec::new(), Vec::new())
    }
    let c: Vec<String> = vec!["?".to_string(); in_.len() - 1];
    let mut conditions = Vec::new();
    conditions.push(format!("{} in ({})", in_[0], c.join(", ")));
    let mut params = Vec::new();
    for i in 1..in_.len() {
        params.push(in_[i].to_string());
    }
    (conditions, params)
}

pub fn lesser_builder(lesser: &[&str]) -> (Vec<String>, Vec<String>) {
    if lesser.len() % 2 != 0 {
        print!("{}\n", "lesser length must be even");
        return (Vec::new(), Vec::new())
    }
    let mut conditions = Vec::new();
    let mut params = Vec::new();
    for i in (0..lesser.len()).step_by(2) {
        let c = format!("{} <= ?", lesser[i]);
        conditions.push(c);
        params.push(lesser[i+1].to_string());
    }
    (conditions, params)
}

pub fn greater_builder(greater: &[&str]) -> (Vec<String>, Vec<String>) {
    if greater.len() % 2 != 0 {
        print!("{}\n", "greater length must be even");
        return (Vec::new(), Vec::new())
    }
    let mut conditions = Vec::new();
    let mut params = Vec::new();
    for i in (0..greater.len()).step_by(2) {
        let c = format!("{} >= ?", greater[i]);
        conditions.push(c);
        params.push(greater[i+1].to_string());
    }
    (conditions, params)
}
