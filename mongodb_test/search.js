// Searching
// like %%
db.products.find({
    name:{
        $regex: /iph/,
        $options: "i"
    }
});

// or like %% -> for product search
db.products.find({
    $or : [
        {
            name:{
                $regex: /iph/,
                $options: "i"
            }
        },{
            description:{
                $regex: /iph/,
                $options: "i"
            }
        }
    ]
});

// Filtering

// Greater than equal
db.products.find({
    $expr: {
        $gte: ["$price", 100000]
    }
});

// Between price
db.products.find({
    price :{
        $gte : 25000,
        $lte : 50000
    }
});

// Price range and other condition
db.products.find({
    price :{
        $gte : 25000,
        $lte : 50000
    },
    category : "peralatan-rumah",
    condition : 1
});

// Sorting
db.products.find({}).sort({
    name: 1,
    category: -1
})


// Full Search with filter and sort
db.products.find({
    $or : [
        {
            name:{
                $regex: /iph/,
                $options: "i"
            }
        },{
            store:{
                $regex: /iph/,
                $options: "i"
            }
        }
    ],
    price :{
        $gte : 25000,
        $lte : 50000
    },
    category : "gadget",
    condition : 1
}).sort({
    name: 1,
    category: -1
});