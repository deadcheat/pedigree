extern crate clap;
extern crate iron;
extern crate params;
extern crate time;

use clap::{Arg, App};
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
    thread::spawn(move || {
                      println!("[LOG]{:?}", obj);
                  });
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

fn main() {
    let matches = App::new("Pedigree")
        .version("1.0")
        .author("deadcheat")
        .about("simple-request logger written in rustlang.")
        .arg(Arg::with_name("host")
                 .short("h")
                 .long("host")
                 .value_name("HOST-NAME")
                 .help("Set hostname of this app. if empty, it'll use localhost.")
                 .required(false))
        .arg(Arg::with_name("port")
                 .short("p")
                 .long("port")
                 .value_name("PORT-NUM")
                 .help("Set portnum of this app. if empty, it'll use 3000.")
                 .required(false))
        .get_matches();
    let host = if matches.is_present("host") {
        matches.value_of("host").unwrap()
    } else {
        "localhost"
    };
    let port = if matches.is_present("port") {
        matches.value_of("port").unwrap()
    } else {
        "3000"
    };
    let mut chain = Chain::new(accept_order);
    chain.link_before(ResponseTime);
    chain.link_after(ResponseTime);
    Iron::new(chain)
        .http(format!("{}:{}", host, port))
        .unwrap();
}
