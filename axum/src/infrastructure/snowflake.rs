use std::sync::Mutex;
use std::time::{SystemTime, UNIX_EPOCH};

const EPOCH: u64 = 1288834974657;
const WORKER_ID_BITS: u64 = 5;
const DATACENTER_ID_BITS: u64 = 5;
const SEQUENCE_BITS: u64 = 12;

const MAX_WORKER_ID: u64 = (1 << WORKER_ID_BITS) - 1;
const MAX_DATACENTER_ID: u64 = (1 << DATACENTER_ID_BITS) - 1;

const WORKER_ID_SHIFT: u64 = SEQUENCE_BITS;
const DATACENTER_ID_SHIFT: u64 = SEQUENCE_BITS + WORKER_ID_BITS;
const TIMESTAMP_SHIFT: u64 = SEQUENCE_BITS + WORKER_ID_BITS + DATACENTER_ID_BITS;

const SEQUENCE_MASK: u64 = (1 << SEQUENCE_BITS) - 1;

pub struct Snowflake {
    worker_id: u64,
    datacenter_id: u64,
    sequence: u64,
    last_timestamp: u64,
    mutex: Mutex<()>,
}

impl Snowflake {
    pub fn new(worker_id: u64, datacenter_id: u64) -> Self {
        if worker_id > MAX_WORKER_ID {
            panic!("worker_id can't be greater than {}", MAX_WORKER_ID);
        }
        if datacenter_id > MAX_DATACENTER_ID {
            panic!("datacenter_id can't be greater than {}", MAX_DATACENTER_ID);
        }

        Snowflake {
            worker_id,
            datacenter_id,
            sequence: 0,
            last_timestamp: 0,
            mutex: Mutex::new(()),
        }
    }

    pub fn next_id(&mut self) -> u64 {
        let _lock = self.mutex.lock().unwrap();
        let mut timestamp = current_timestamp();

        if timestamp < self.last_timestamp {
            panic!("Clock moved backwards. Refusing to generate id");
        }

        if self.last_timestamp == timestamp {
            self.sequence = (self.sequence + 1) & SEQUENCE_MASK;
            if self.sequence == 0 {
                timestamp = self.wait_for_next_millis(self.last_timestamp);
            }
        } else {
            self.sequence = 0;
        }

        self.last_timestamp = timestamp;

        ((timestamp - EPOCH) << TIMESTAMP_SHIFT)
            | (self.datacenter_id << DATACENTER_ID_SHIFT)
            | (self.worker_id << WORKER_ID_SHIFT)
            | self.sequence
    }

    fn wait_for_next_millis(&self, last_timestamp: u64) -> u64 {
        let mut timestamp = current_timestamp();
        while timestamp <= last_timestamp {
            timestamp = current_timestamp();
        }
        timestamp
    }
}

fn current_timestamp() -> u64 {
    let start = SystemTime::now();
    let since_the_epoch = start
        .duration_since(UNIX_EPOCH)
        .expect("Time went backwards");
    since_the_epoch.as_millis() as u64
}

fn main() {
    let mut snowflake = Snowflake::new(1, 1);
    let id = snowflake.next_id();
    println!("Generated ID: {}", id);
}
