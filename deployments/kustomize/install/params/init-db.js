function env(name, fallback) {
    const value = process.env[name];
    return value === undefined || value === "" ? fallback : value;
}

const mongoHost = env("AMBULANCE_MANAGEMENT_API_MONGODB_HOST", env("AMBULANCE_API_MONGODB_HOST", "mongodb"));
const mongoPort = env("AMBULANCE_MANAGEMENT_API_MONGODB_PORT", env("AMBULANCE_API_MONGODB_PORT", "27017"));

const mongoUser = env("AMBULANCE_MANAGEMENT_API_MONGODB_USERNAME", env("AMBULANCE_API_MONGODB_USERNAME", ""));
const mongoPassword = env("AMBULANCE_MANAGEMENT_API_MONGODB_PASSWORD", env("AMBULANCE_API_MONGODB_PASSWORD", ""));

const databaseName = env("AMBULANCE_MANAGEMENT_API_MONGODB_DATABASE", env("AMBULANCE_API_MONGODB_DATABASE", "ambulance_management"));
const vehiclesCollectionName = env("AMBULANCE_MANAGEMENT_API_MONGODB_VEHICLES_COLLECTION", "vehicles");
const countersCollectionName = env("AMBULANCE_MANAGEMENT_API_MONGODB_COUNTERS_COLLECTION", "counters");

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5", 10) || 5;

function mongoUri() {
    if (!mongoUser) {
        return `mongodb://${mongoHost}:${mongoPort}`;
    }
    return `mongodb://${encodeURIComponent(mongoUser)}:${encodeURIComponent(mongoPassword)}@${mongoHost}:${mongoPort}/?authSource=admin`;
}

let connection;
while (true) {
    try {
        connection = Mongo(mongoUri());
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`);
        sleep(retrySeconds * 1000);
    }
}

const db = connection.getDB(databaseName);
const existingCollections = db.getCollectionNames();

if (!existingCollections.includes(vehiclesCollectionName)) {
    db.createCollection(vehiclesCollectionName);
}

if (!existingCollections.includes(countersCollectionName)) {
    db.createCollection(countersCollectionName);
}

db[vehiclesCollectionName].createIndex({ vehicleId: 1 }, { unique: true });
db[vehiclesCollectionName].createIndex({ callSign: 1 }, { unique: true });
db[vehiclesCollectionName].createIndex({ plateNumber: 1 }, { unique: true });

const sampleVehicles = [
    {
        vehicleId: 1,
        callSign: "AMB-101",
        vehicleType: "Type B ambulance",
        plateNumber: "BA-101XY",
        station: "Bratislava Center",
        assignedCrew: "Novak / Simko",
        status: "AVAILABLE",
        mileageKm: 48210,
        lastServiceDate: "2026-03-15",
        notes: "Primary city response unit."
    },
    {
        vehicleId: 2,
        callSign: "AMB-204",
        vehicleType: "Type C ambulance",
        plateNumber: "TT-204LM",
        station: "Trnava North",
        assignedCrew: "Kovac / Balaz",
        status: "ON_MISSION",
        mileageKm: 73125,
        lastServiceDate: "2026-02-10",
        notes: "Advanced life support vehicle."
    },
    {
        vehicleId: 3,
        callSign: "AMB-309",
        vehicleType: "Type B ambulance",
        plateNumber: "NR-309AB",
        station: "Nitra South",
        assignedCrew: "Unassigned",
        status: "IN_SERVICE",
        mileageKm: 61540,
        lastServiceDate: "2026-04-01",
        notes: "Scheduled cleaning and equipment restock."
    },
    {
        vehicleId: 4,
        callSign: "AMB-412",
        vehicleType: "Type A ambulance",
        plateNumber: "ZA-412CD",
        station: "Zilina East",
        assignedCrew: "Mraz / Halas",
        status: "OUT_OF_SERVICE",
        mileageKm: 91880,
        lastServiceDate: "2026-01-20",
        notes: "Waiting for transmission repair."
    }
];

if (db[vehiclesCollectionName].countDocuments({}) === 0) {
    const insertResult = db[vehiclesCollectionName].insertMany(sampleVehicles);
    if (insertResult.writeError) {
        printjson(insertResult);
        throw new Error(`Error when writing vehicle seed data: ${insertResult.errmsg}`);
    }
    print(`Inserted ${sampleVehicles.length} sample vehicles into '${vehiclesCollectionName}'`);
} else {
    print(`Collection '${vehiclesCollectionName}' already contains data, skipping vehicle seed`);
}

const maxVehicleIdDoc = db[vehiclesCollectionName]
    .find({}, { vehicleId: 1, _id: 0 })
    .sort({ vehicleId: -1 })
    .limit(1)
    .toArray()[0];
const nextSeq = maxVehicleIdDoc ? maxVehicleIdDoc.vehicleId : 0;

db[countersCollectionName].updateOne(
    { _id: "vehicles" },
    { $set: { seq: nextSeq } },
    { upsert: true }
);

print(`Counter '${countersCollectionName}/vehicles' set to ${nextSeq}`);

process.exit(0);
