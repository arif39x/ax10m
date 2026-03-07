module MathOracle

using HTTP
using JSON3
using StructTypes
using LinearAlgebra

mutable struct ComputeRequest
    x::Vector{Float32}
    y::Vector{Float32}
    z::Vector{Float32}
    field_x::Vector{Float32}
    field_y::Vector{Float32}
    field_z::Vector{Float32}
    field_amplitude::Vector{Float32}
    field_sigma::Vector{Float32}
    mass::Vector{Float32}
    
    ComputeRequest() = new([], [], [], [], [], [], [], [], [])
end

mutable struct ComputeResponse
    ax::Vector{Float32}
    ay::Vector{Float32}
    az::Vector{Float32}
end

StructTypes.StructType(::Type{ComputeRequest}) = StructTypes.Mutable()
StructTypes.StructType(::Type{ComputeResponse}) = StructTypes.Struct()

function compute_acceleration(req::ComputeRequest)::ComputeResponse
    n_entities = length(req.x)
    n_fields = length(req.field_x)
    
    ax = zeros(Float32, n_entities)
    ay = zeros(Float32, n_entities)
    az = zeros(Float32, n_entities)
    
    for i in 1:n_entities
        r_i = [req.x[i], req.y[i], req.z[i]]
        net_a = [0.0f0, 0.0f0, -9.8f0]
        m_i = (i <= length(req.mass) && req.mass[i] != 0.0f0) ? req.mass[i] : 1.0f0
        
        for j in 1:n_fields
            r_j = [req.field_x[j], req.field_y[j], req.field_z[j]]
            A_j = req.field_amplitude[j]
            sigma_j = req.field_sigma[j]
            
            diff_vec = r_i .- r_j
            dist_sq = dot(diff_vec, diff_vec)
            
            val = A_j * exp(-dist_sq / (2.0f0 * sigma_j^2))
            grad = - (diff_vec ./ (sigma_j^2)) .* val
            
            net_a .-= (1.0f0 / m_i) .* grad
        end
        
        ax[i] = net_a[1]
        ay[i] = net_a[2]
        az[i] = net_a[3]
    end
    
    return ComputeResponse(ax, ay, az)
end

function handle_request(req::HTTP.Request)
    try
        comp_req = JSON3.read(req.body, ComputeRequest)
        comp_resp = compute_acceleration(comp_req)
        out_bytes = JSON3.write(comp_resp)
        return HTTP.Response(200, ["Content-Type" => "application/json"], out_bytes)
    catch e
        println("Error processing request: ", e)
        return HTTP.Response(400, "Bad Request")
    end
end

function main()
    println("Julia Math Oracle (JSON) listening on 0.0.0.0:50051")
    HTTP.serve(handle_request, "0.0.0.0", 50051)
end

end

MathOracle.main()
