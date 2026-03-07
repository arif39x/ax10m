struct State {
    x: f32,
    y: f32,
    z: f32,
    padding: f32,
}

@group(0) @binding(0)
var<uniform> state: State;

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) uv: vec2<f32>,
};

@vertex
fn vs_main(@builtin(vertex_index) in_vertex_index: u32) -> VertexOutput {
    var out: VertexOutput;
    let x = f32((in_vertex_index << 1) & 2u);
    let y = f32(in_vertex_index & 2u);
    out.clip_position = vec4<f32>(x * 2.0 - 1.0, 1.0 - y * 2.0, 0.0, 1.0);
    out.uv = vec2<f32>(x, y);
    return out;
}

fn map(p: vec3<f32>) -> f32 {
    let sphere_pos = vec3<f32>(state.x, state.y, state.z);
    return length(p - sphere_pos) - 10.0;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    let uv = in.uv * 2.0 - 1.0;
    
    let ro = vec3<f32>(0.0, 0.0, 200.0);
    let rd = normalize(vec3<f32>(uv.x, uv.y, -1.0));
    
    var t = 0.0;
    for (var i = 0; i < 64; i = i + 1) {
        let p = ro + rd * t;
        let d = map(p);
        if (d < 0.001) {
            let col = vec3<f32>(0.5, 0.8, 1.0) * (1.0 - f32(i)/64.0);
            return vec4<f32>(col, 1.0);
        }
        t = t + d;
        if (t > 400.0) {
            break;
        }
    }
    
    return vec4<f32>(0.0, 0.0, 0.0, 1.0);
}
