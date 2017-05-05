extern crate iron;
extern crate time;
extern crate fruently;
extern crate params;

use fruently::fluent::Fluent;
use std::collections::HashMap;
use fruently::forwardable::JsonForwardable;
use iron::prelude::*;
use iron::{BeforeMiddleware, AfterMiddleware, typemap};
use time::precise_time_ns;
use std::io::Read;
use params::Params;

struct ResponseTime;

impl typemap::Key for ResponseTime { type Value = u64; }

impl BeforeMiddleware for ResponseTime {
    fn before(&self, req: &mut Request) -> IronResult<()> {
        req.extensions.insert::<ResponseTime>(precise_time_ns());
        Ok(())
    }
}

impl AfterMiddleware for ResponseTime {
    fn after(&self, req: &mut Request, res: Response) -> IronResult<Response> {
        let delta = precise_time_ns() - *req.extensions.get::<ResponseTime>().unwrap();
        println!("Request took: {} ms", (delta as f64) / 1000000.0);
        Ok(res)
    }
}

fn hello_world(req: &mut Request) -> IronResult<Response> {
    let mut obj: HashMap<String, String> = HashMap::new();

    obj.insert("Method".to_string(), req.method.to_string());
    for header in req.headers.iter() {
        obj.insert(header.name().to_string(), header.value_string());
    }

    let map = req.get_ref::<Params>().unwrap();
    obj.insert("Params".to_string(), format!("{:?}", map));

    // let mut payload = String::new();
    // req.body.read_to_string(&mut payload).unwrap();
    // obj.insert("Body".to_string(), payload);
    let fruently = Fluent::new("127.0.0.1:24224", "test");
    match fruently.post(&obj) {
        Err(e) => println!("{:?}", e),
        Ok(_) => println!("{:?}", obj),
    }
    Ok(Response::with((iron::status::Ok, "Hello World")))
}

fn main() {
    let mut chain = Chain::new(hello_world);
    chain.link_before(ResponseTime);
    chain.link_after(ResponseTime);
    Iron::new(chain).http("localhost:3000").unwrap();
}
