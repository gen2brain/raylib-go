#version 330

#define MAX_LIGHTS              4
#define LIGHT_DIRECTIONAL       0
#define LIGHT_POINT             1
#define LIGHT_SPOT 2
#define PI 3.14159265358979323846

struct Light{
    int enabled;
    int type;
    float enargy;
    float cutOff ;
    float outerCutOff;
    float constant;
    float linear ;
    float quadratic ;
    float shiny ;
    float specularStr;
    vec3 position;
    vec3 direction;
    vec3 lightColor;
};

// Input vertex attributes (from vertex shader)
in vec3 fragPosition;
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragNormal;
in vec4 shadowPos;
in mat3 TBN;

// Output fragment color
out vec4 finalColor;
// mask
uniform sampler2D mask;
uniform int frame;
// Input uniform values
uniform int numOfLights = 4;
uniform sampler2D albedoMap;
uniform sampler2D mraMap;
uniform sampler2D normalMap;
uniform sampler2D emissiveMap; // r: Hight g:emissive
// Input uniform values
uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D flashlight;

uniform vec4 colDiffuse;

uniform vec2 tiling = vec2(0.5);
uniform vec2 offset ;
uniform vec2 tilingFlashlight = vec2(0.5);
uniform vec2 offsetFlashlight ;

uniform int useTexAlbedo =1;
uniform int useTexNormal = 0;
uniform int useTexMRA =1;
uniform int useTexEmissive =1;

uniform vec4  albedoColor ;
uniform vec4  emissiveColor ;
uniform float normalValue =0.5;
uniform float metallicValue =0.4;
uniform float roughnessValue =0;
uniform float aoValue =0.8;
uniform float emissivePower ;

// Input lighting values
uniform Light lights[MAX_LIGHTS]; 
uniform vec3 viewPos;

uniform vec3 ambientColor = vec3(1,1,1);
uniform float ambientStrength = 0.2;
uniform float ambient = 0.03;
uniform float fogDensity;

vec3 CalcDirLight(Light light,vec3 normal,vec3 viewDir,vec3 albedo,vec3 baseRefl,float roughness,float metallic);
vec3 CalcPointLight(Light light,vec3 normal,vec3 viewDir,vec3 albedo,vec3 baseRefl,float roughness,float metallic);
vec3 CalcSpotLight(Light light,vec3 normal,vec3 viewDir,vec3 albedo,vec3 baseRefl,float roughness,float metallic);

// Reflectivity in range 0.0 to 1.0
// NOTE: Reflectivity is increased when surface view at larger angle
vec3 SchlickFresnel(float hDotV,vec3 refl)
{
    return refl + (1.0 - refl)*pow(1.0 - hDotV, 5.0);
}

float GgxDistribution(float nDotH,float roughness)
{
    float a = roughness * roughness * roughness * roughness;
    float d = nDotH * nDotH * (a - 1.0) + 1.0;
    d = PI * d * d;
    return a / max(d,0.0000001);
}

float GeomSmith(float nDotV,float nDotL,float roughness)
{
    float r = roughness + 1.0;
    float k = r*r / 8.0;
    float ik = 1.0 - k;
    float ggx1 = nDotV/(nDotV*ik + k);
    float ggx2 = nDotL/(nDotL*ik + k);
    return ggx1*ggx2;
}

vec3 ComputePBR()
{
     vec3 albedo = texture(albedoMap,vec2(fragTexCoord.x*tiling.x + offset.x, fragTexCoord.y*tiling.y + offset.y)).rgb;
     albedo = vec3(albedoColor.x*albedo.x, albedoColor.y*albedo.y, albedoColor.z*albedo.z);
    
    float metallic = clamp(metallicValue, 0.0, 1.0);
    float roughness = clamp(roughnessValue, 0.0, 1.0);
    float ao = clamp(aoValue, 0.0, 1.0);
    
    if (useTexMRA == 1)
    {
        vec4 mra = texture(mraMap, vec2(fragTexCoord.x*tiling.x + offset.x, fragTexCoord.y*tiling.y + offset.y))*useTexMRA;
        metallic = clamp(mra.r + metallicValue, 0.04, 1.0);
        roughness = clamp(mra.g + roughnessValue, 0.04, 1.0);
        ao = (mra.b + aoValue)*0.5;
    }

    vec3 N = normalize(fragNormal);
    if (useTexNormal == 1)
    {
        N = texture(normalMap, vec2(fragTexCoord.x*tiling.x + offset.y, fragTexCoord.y*tiling.y + offset.y)).rgb;
        N = normalize(N*2.0 - 1.0);
        N = normalize(N*TBN);
    }

    vec3 V = normalize(viewPos - fragPosition);

    vec3 emissive = vec3(0);
    emissive = (texture(emissiveMap, vec2(fragTexCoord.x*tiling.x+offset.x, fragTexCoord.y*tiling.y+offset.y)).rgb).g * emissiveColor.rgb*emissivePower * useTexEmissive;

    // return N;//vec3(metallic,metallic,metallic);
    // if dia-electric use base reflectivity of 0.04 otherwise ut is a metal use albedo as base reflectivity
    vec3 baseRefl = mix(vec3(0.04), albedo.rgb, metallic);
    vec3 lightAccum = vec3(0.0);  // Acumulate lighting lum

     vec3 norm=N;
     vec3 viewDir=V;
     vec3 result = vec3(0.0);
   
       for (int i = 0; i < MAX_LIGHTS; i++){
       
        if(lights[i].enabled == 1){

            if(lights[i].type == LIGHT_DIRECTIONAL){
                result += CalcDirLight(lights[i],norm,viewDir,albedo,baseRefl,roughness,metallic);  
            }

            if(lights[i].type == LIGHT_POINT){
                result += CalcPointLight(lights[i],norm,viewDir,albedo,baseRefl,roughness,metallic);
            }

            if(lights[i].type == LIGHT_SPOT){
                 result += CalcSpotLight(lights[i],norm,viewDir,albedo,baseRefl,roughness,metallic);   
            }

        }

    } 
       
    vec3 ambientFinal = (ambientColor + albedo)*ambient*0.5;
    
    return ambientFinal+result*ao + emissive;

}


void main()
{
    vec3 color = ComputePBR();

    // HDR tonemapping
    color = pow(color, color + vec3(1.0));
    
    // // Gamma correction
   //  color = pow(color, vec3(1.0/2.5));

    // Fog calculation
    float dist = length(viewPos - fragPosition);

    // these could be parameters...
    const vec4 fogColor = vec4(0.5, 0.5, 0.5, 1.0);


    // Exponential fog
    float fogFactor = 1.0/exp((dist*fogDensity)*(dist*fogDensity));


    fogFactor = clamp(fogFactor, 0.0, 1.0);

    finalColor = mix(fogColor,vec4(color,1.0), fogFactor);
}



vec3 CalcDirLight(Light light,vec3 normal,vec3 viewDir,vec3 albedo,vec3 baseRefl,float roughness,float metallic)
{

    vec3 L = normalize(-light.direction);
    float diff=max(dot(normal,L),0.0);
    vec3 diffuse=light.lightColor*diff*vec3(texture(texture0,fragTexCoord));      
    vec3 H = normalize(diffuse + L);

           // Cook-Torrance BRDF distribution function
    float nDotV = max(dot(normal,viewDir), 0.0000001);
    float nDotL = max(dot(normal,L), 0.0000001);
    float hDotV = max(dot(H,viewDir), 0.0);
    float nDotH = max(dot(normal,H), 0.0);
    float D = GgxDistribution(nDotH, roughness);    // Larger the more micro-facets aligned to H
    float G = GeomSmith(nDotV, nDotL, roughness);   // Smaller the more micro-facets shadow
    vec3 F = SchlickFresnel(hDotV, baseRefl);       // Fresnel proportion of specular reflectance  

    vec3 spec = (D*G*F)/(4.0*nDotV*nDotL);
        
     // Difuse and spec light can't be above 1.0
    // kD = 1.0 - kS  diffuse component is equal 1.0 - spec comonent
    vec3 kD = vec3(1.0) - F;
    // Mult kD by the inverse of metallnes, only non-metals should have diffuse light
    kD *= 1.0 - metallic;
    // Angle of light has impact on result
    return ((kD*albedo.rgb/PI + spec)*nDotL)*light.enabled;

}

vec3 CalcPointLight(Light light,vec3 normal,vec3 viewDir,vec3 albedo,vec3 baseRefl,float roughness,float metallic)
{
    
        vec3 L = normalize(light.position - fragPosition);      
        vec3 H = normalize(viewDir + L);                                
        float distance=length(light.position-fragPosition);
        float attenuation=light.enargy/(light.constant+light.linear*distance+light.quadratic*(distance*distance)); 
        vec3 radiance = light.lightColor.rgb*light.enargy*attenuation; 

        // Cook-Torrance BRDF distribution function
        float nDotV = max(dot(normal,viewDir), 0.0000001);
        float nDotL = max(dot(normal,L), 0.0000001);
        float hDotV = max(dot(H,viewDir), 0.0);
        float nDotH = max(dot(normal,H), 0.0);
        float D = GgxDistribution(nDotH, roughness);    // Larger the more micro-facets aligned to H
        float G = GeomSmith(nDotV, nDotL, roughness);   // Smaller the more micro-facets shadow
        vec3 F = SchlickFresnel(hDotV, baseRefl);       // Fresnel proportion of specular reflectance

        vec3 spec = (D*G*F)/(4.0*nDotV*nDotL);
        
        // Difuse and spec light can't be above 1.0
        // kD = 1.0 - kS  diffuse component is equal 1.0 - spec comonent
        vec3 kD = vec3(1.0) - F;
        // Mult kD by the inverse of metallnes, only non-metals should have diffuse light
        kD *= 1.0 - metallic;
        // Angle of light has impact on result
    return ((kD*albedo.rgb/PI + spec)*radiance*nDotL)*light.enabled ;
}
    
vec3 CalcSpotLight(Light light,vec3 normal,vec3 viewDir,vec3 albedo,vec3 baseRefl,float roughness,float metallic){

        vec3 L = normalize(light.position - fragPosition);

        float theta=dot(L,normalize(-light.direction));
        float epsilon=cos(radians(light.cutOff))-cos(radians(light.outerCutOff));
        float intensity=smoothstep(0.0,1.0,(theta-cos(radians(light.outerCutOff)))/epsilon);//clamp((theta-cos(radians(light.outerCutOff)))/epsilon,0.0,1.0);
        intensity*= length(vec3(texture(flashlight,vec2(fragTexCoord.x*tilingFlashlight.x + offsetFlashlight.y, fragTexCoord.y*tilingFlashlight.y + offsetFlashlight.y)).rgb));

        float diff=max(dot(normal,L),0.0);      
        vec3 H = light.lightColor*diff*vec3(texture(texture0,fragTexCoord));
        float distance=length(light.position-fragPosition);
        float attenuation=light.enargy/(light.constant+light.linear*distance+light.quadratic*(distance*distance)); 
        vec3 radiance = light.lightColor.rgb*light.enargy*attenuation; 
        H*=intensity;    

        // Cook-Torrance BRDF distribution function
        float nDotV = max(dot(normal,viewDir), 0.0000001);
        float nDotL = max(dot(normal,L), 0.0000001);
        float hDotV = max(dot(H,viewDir), 0.0);
        float nDotH = max(dot(normal,H), 0.0);
        float D = GgxDistribution(nDotH, roughness);    // Larger the more micro-facets aligned to H
        float G = GeomSmith(nDotV, nDotL, roughness);   // Smaller the more micro-facets shadow
        vec3 F = SchlickFresnel(hDotV, baseRefl);       // Fresnel proportion of specular reflectance

       vec3 spec = (D*G*F)/(4.0*nDotV*nDotL);
       spec*=intensity; 
        
        
        // Difuse and spec light can't be above 1.0
        // kD = 1.0 - kS  diffuse component is equal 1.0 - spec comonent
        vec3 kD = vec3(1.0) - F;
        // Mult kD by the inverse of metallnes, only non-metals should have diffuse light
        kD *= 1.0 - metallic;
        // Angle of light has impact on result
        return ((kD*albedo.rgb/PI + spec)*radiance*nDotL)*light.enabled;
}

