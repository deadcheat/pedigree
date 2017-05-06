extern crate fruently;
extern crate iron;
extern crate params;
extern crate time;

use fruently::fluent::Fluent;
use fruently::forwardable::JsonForwardable;
use iron::prelude::*;
use iron::status;
use iron::{BeforeMiddleware, AfterMiddleware, typemap};
use params::Params;
use time::precise_time_ns;
use std::collections::HashMap;
use std::thread;


struct ResponseTime;

impl typemap::Key for ResponseTime {
    type Value = u64;
}

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

fn accept_order(req: &mut Request) -> IronResult<Response> {
    let obj: HashMap<String, String> = pack_info(req);
    thread::spawn(move || { logging_request(obj); });
    Ok(Response::with(status::Created))
}

fn pack_info(req: &mut Request) -> HashMap<String, String> {
    let mut obj: HashMap<String, String> = HashMap::new();

    obj.insert("Method".to_string(), req.method.to_string());
    for header in req.headers.iter() {
        obj.insert(header.name().to_string(), header.value_string());
    }

    let map = req.get_ref::<Params>().unwrap();
    obj.insert("Params".to_string(), format!("{:?}", map));

    return obj;
}

fn logging_request(obj: HashMap<String, String>) {

    let fruently = Fluent::new("127.0.0.1:24224", "test");
    match fruently.post(&obj) {
        Err(e) => println!("[ERR]{:?}", e),
        Ok(_) => println!("[SUC]{:?}", obj),
    }
}

fn main() {
    let mut chain = Chain::new(accept_order);
    chain.link_before(ResponseTime);
    chain.link_after(ResponseTime);
    Iron::new(chain).http("localhost:3000").unwrap();
}
